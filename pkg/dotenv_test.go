package dotenv

import (
	"os"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLoadEnvFile(t *testing.T) {

	tests := []struct {
		name        string
		content     string
		expected    map[string]string
		expectedErr bool
	}{
		{
			name:    "Valid Case",
			content: "KEY1=value1\nKEY2=value2\nKEY3=value3\n",
			expected: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
				"KEY3": "value3",
			},
			expectedErr: false,
		},
		{
			name:    "Case with Comments",
			content: "# Comment\nKEY21=value21\n# Another comment\nKEY22=value22\n",
			expected: map[string]string{
				"KEY21": "value21",
				"KEY22": "value22",
			},
			expectedErr: false,
		},
		{
			name:    "Comment in Same Line with Key-Value Pair",
			content: "KEY31=value31 # Comment\nKEY32=value32\n",
			expected: map[string]string{
				"KEY31": "value31",
				"KEY32": "value32",
			},
			expectedErr: false,
		},
		{
			name:    "Case Using : Separator",
			content: "KEY41:value41\nKEY42:value42\n",
			expected: map[string]string{
				"KEY41": "value41",
				"KEY42": "value42",
			},
			expectedErr: false,
		},
		{
			name:    "Case Export Keyword",
			content: "export KEY51=value51\nexport KEY52=value52\n",
			expected: map[string]string{
				"KEY51": "value51",
				"KEY52": "value52",
			},
			expectedErr: false,
		},
		{
			name:        "Case Key Length = 0",
			content:     "=value61\nKEY62=value62\n",
			expectedErr: true,
		},
		{
			name:        "Case Value Length = 0",
			content:     "KEY71=\nKEY72=value72\n",
			expectedErr: true,
		},
		{
			name:    "Case Using unvalid Separator",
			content: "KEY71#value71\nKEY72#value72\n",
			expected: map[string]string{
				"KEY71": "value71",
				"KEY72": "value72",
			},
			expectedErr: true,
		},
	}

	for i, test := range tests {
		file, err := os.Create(t.TempDir() + "/" + ".env" + strconv.Itoa(i))
		if err != nil {
			assert.NoError(t, err)
		}
		_, err = file.WriteString(test.content)
		if err != nil {
			assert.NoError(t, err)
		}
		file.Close()
		t.Run(test.name, func(t *testing.T) {
			err = Load(file.Name())
			if test.expectedErr {
				assert.Error(t, err)
				return
			}
			for k, v := range test.expected {
				assert.Equal(t, v, os.Getenv(k))
			}

		})
	}

	t.Run("case with default file given", func(t *testing.T) {
		file, err := os.Create(t.TempDir() + "/" + ".env")
		if err != nil {
			assert.NoError(t, err)
		}
		_, err = file.WriteString(tests[0].content)
		if err != nil {
			assert.NoError(t, err)
		}
		file.Close()
		err = Load()
		if tests[0].expectedErr {
			assert.Error(t, err)
		}
		for k, v := range tests[0].expected {
			assert.Equal(t, v, os.Getenv(k))
		}

	})
}

func TestUnmarshal(t *testing.T) {
	tests := []struct {
		name        string
		content     string
		expected    map[string]string
		expectedErr bool
	}{
		{
			name:    "Valid Case",
			content: "KEY1=value1\nKEY2=value2\nKEY3=value3\n",
			expected: map[string]string{
				"KEY1": "value1",
				"KEY2": "value2",
				"KEY3": "value3",
			},
			expectedErr: false,
		},
		{
			name:    "Case with Comments",
			content: "# Comment\nKEY21=value21\n# Another comment\nKEY22=value22\n",
			expected: map[string]string{
				"KEY21": "value21",
				"KEY22": "value22",
			},
			expectedErr: false,
		},
		{
			name:    "Comment in Same Line with Key-Value Pair",
			content: "KEY31=value31 # Comment\nKEY32=value32\n",
			expected: map[string]string{
				"KEY31": "value31",
				"KEY32": "value32",
			},
			expectedErr: false,
		},
		{
			name:    "Case Using : Separator",
			content: "KEY41:value41\nKEY42:value42\n",
			expected: map[string]string{
				"KEY41": "value41",
				"KEY42": "value42",
			},
			expectedErr: false,
		},
		{
			name:    "Case Export Keyword",
			content: "export KEY51=value51\nexport KEY52=value52\n",
			expected: map[string]string{
				"KEY51": "value51",
				"KEY52": "value52",
			},
			expectedErr: false,
		},
		{
			name:        "Case Key Length = 0",
			content:     "=value61\nKEY62=value62\n",
			expectedErr: true,
		},
		{
			name:        "Case Value Length = 0",
			content:     "KEY71=\nKEY72=value72\n",
			expectedErr: true,
		},
		{
			name:    "Case Using unvalid Separator",
			content: "KEY71#value71\nKEY72#value72\n",
			expected: map[string]string{
				"KEY71": "value71",
				"KEY72": "value72",
			},
			expectedErr: true,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := Unmarshal(test.content)
			if test.expectedErr {
				assert.Error(t, err)
				return
			}

			for k, v := range test.expected {
				assert.Equal(t, v, result[k])
			}
		})
	}
}
