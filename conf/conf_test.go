package conf

import (
	"fmt"
	A "github.com/stretchr/testify/assert"
	"os"
	"path/filepath"
	"testing"
)

var (
	mockCfgDir   = "/tmp"
	mockCfgName  = "cfg_test"
	mockCfgPath  = filepath.Join(mockCfgDir, fmt.Sprintf("%s.%s", mockCfgName, "yaml"))
	mockTemplate = `
agent:
  # encode protocol. optional: json, msgpack, default: msgpack
  encode: {{encode|default:"msgpack"}}
    
  # whether using gzip to compress the raw data. default: false
  # if gzip is true, using lz4 to compress data.
  gzip: {{gzip|default:false}}

  topics:
    {% if topics %}
    {% for topic in topics %}
    - stream: {{topic.name}}
      subjects:
        {% for subject in topic.subjects %}
        - {{subject}}
        {% endfor %}
    {% endfor %}
    {% endif %}
`
)

func testWrite(t *testing.T) {
	assert := A.New(t)

	var (
		expect = `
agent:
  # encode protocol. optional: json, msgpack, default: msgpack
  encode: msgpack

  # whether using gzip to compress the raw data. default: false
  # if gzip is true, using lz4 to compress data.
  gzip: false

  topics:
    - stream: a
      subjects:
        - 0.0.0.0
    - stream: b
      subjects:
        - 1.1.1.1
`
	)

	c := New(
		WithConfigType("yaml"),
		WithAutomaticEnv(),
		WithAllowEmptyEnv(),
		WithTemplate(mockTemplate),
		WithOptionConfigPath([]string{mockCfgPath}),
		WithWriteTo(mockCfgPath),
		WithDefaultVal(map[string]interface{}{
			"topics": []map[string]interface{}{
				{
					"name": "a",
					"subjects": []string{
						"0.0.0.0",
					},
				},
				{
					"name": "b",
					"subjects": []string{
						"1.1.1.1",
					},
				},
			},
		}),
	)

	err := c.Write()
	assert.Nil(err)

	buf, err := os.ReadFile(mockCfgPath)
	assert.Nil(err)
	assert.Equal(string(buf), expect)
}

func testRead(t *testing.T) {
	assert := A.New(t)

	c := New(
		WithConfigType("yaml"),
		WithAutomaticEnv(),
		WithTemplate(mockTemplate),
		WithOptionConfigPath([]string{mockCfgPath}),
		WithWriteTo(mockCfgPath),
		WithDefaultVal(map[string]interface{}{
			"topics": []map[string]interface{}{
				{
					"name": "a",
					"subjects": []string{
						"0.0.0.0",
					},
				},
				{
					"name": "b",
					"subjects": []string{
						"1.1.1.1",
					},
				},
			},
		}),
	)

	err := c.Read()
	assert.Nil(err)

	assert.Equal(c.GetString("agent.encode"), "msgpack")
	topics := c.GetMapSlice("agent.topics")
	for _, topic := range topics {
		switch topic["name"] {
		case "a":
			assert.Equal("0.0.0.0", topic["pattern"])
		case "b":
			assert.Equal("1.1.1.1", topic["pattern"])
		}
	}

}

func TestCfg(t *testing.T) {
	t.Run("TestWrite", testWrite)
	t.Run("TestRead", testRead)
}
