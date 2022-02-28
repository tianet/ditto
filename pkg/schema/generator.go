package schema

import (
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/lucasjones/reggen"
	"github.com/patrickmn/go-cache"
	"gonum.org/v1/gonum/stat/distuv"
)

func (f Field) generateUUID() (string, error) {
	return uuid.New().String(), nil
}

func (f Field) generateString() (string, error) {
	if f.Value != nil {
		return (f.Value).(string), nil
	}

	if getRandomProbability() < f.Params.Empty {
		return "", nil
	}

	if f.Params.Regex != "" {
		return reggen.Generate(f.Params.Regex, 64)
	}

	return "", nil
}

func (f Field) generateInt() (int, error) {
	if f.Value != nil {
		return (f.Value).(int), nil
	}

	if f.Params.Distribution != "" {

		switch f.Params.Distribution {

		case PARAMS_DISTRIBUTION_NORMAL:
			return int(distuv.Normal{Mu: f.Params.Mu, Sigma: f.Params.Sigma}.Rand()), nil

		case PARAMS_DISTRIBUTION_UNIFORM:
			return int(distuv.Uniform{Min: f.Params.Min, Max: f.Params.Max}.Rand()), nil

		case PARAMS_DISTRIBUTION_RANDOM:
			if f.Params.Max != 0 {
				return rand.Intn(int(f.Params.Max)), nil
			}
			return rand.Int(), nil

		default:
			return 0, fmt.Errorf("Field %s doesn't have a valid distribuition value", f.Name)
		}
	}

	if f.Params.Incremental {
		if value, ok := C.Get(f.uuid); ok {
			var step int
			if f.Params.Sigma != 0 {
				step = int(distuv.Normal{Mu: f.Params.Step.(float64), Sigma: f.Params.Sigma}.Rand())
			} else {
				step = int((f.Params.Step).(float64))
			}
			new_value := value.(int) + step
			C.Set(f.uuid, new_value, cache.NoExpiration)
			return new_value, nil
		} else {
			if f.Params.Start == nil {
				return 0, fmt.Errorf("Field %s doesn't have a starting point", f.Name)
			}
			start_int := int((f.Params.Start).(float64))
			C.Set(f.uuid, start_int, cache.NoExpiration)
			return start_int, nil
		}
	}

	return 0, fmt.Errorf("Field %s doesn't have a valid configuration", f.Name)
}

func roundFloat(value float64, decimals int) (float64, error) {
	format_string := fmt.Sprintf("%s%s%df", "%", ".", decimals)

	i := fmt.Sprintf(format_string, value)
	f, err := strconv.ParseFloat(i, 64)

	return f, err
}

func (f Field) generateFloat() (float64, error) {
	if f.Value != nil {
		return (f.Value).(float64), nil
	}

	var generated_float float64
	if f.Params.Distribution != "" {

		switch f.Params.Distribution {

		case PARAMS_DISTRIBUTION_NORMAL:
			generated_float = distuv.Normal{Mu: f.Params.Mu, Sigma: f.Params.Sigma}.Rand()

		case PARAMS_DISTRIBUTION_UNIFORM:
			generated_float = distuv.Uniform{Min: f.Params.Min, Max: f.Params.Max}.Rand()

		case PARAMS_DISTRIBUTION_RANDOM:
			generated_float = rand.Float64()
			if f.Params.Scale != 0 {
				generated_float = generated_float * f.Params.Scale
			}

		default:
			return 0, fmt.Errorf("Field %s doesn't have a valid distribuition value", f.Name)
		}

	}

	if f.Params.Incremental {
		if value, ok := C.Get(f.uuid); ok {
			var step float64
			if f.Params.Sigma != 0 {
				step = distuv.Normal{Mu: f.Params.Step.(float64), Sigma: f.Params.Sigma}.Rand()
			} else {
				step = (f.Params.Step).(float64)
			}
			new_value := value.(float64) + step

			C.Set(f.uuid, new_value, cache.NoExpiration)
			generated_float = new_value
		} else {
			if f.Params.Start == nil {
				return 0, fmt.Errorf("Field %s doesn't have a starting point", f.Name)
			}
			start_float := (f.Params.Start).(float64)
			C.Set(f.uuid, start_float, cache.NoExpiration)
			generated_float = start_float
		}
	}

	if generated_float == 0 {
		return 0, fmt.Errorf("Field %s doesn't have a valid configuration", f.Name)
	}

	round := 2
	if f.Params.Round != 0 {
		round = f.Params.Round
	}

	generated_float, err := roundFloat(generated_float, round)
	if err != nil {
		return generated_float, err
	}

	return generated_float, nil
}

func getTimeStep(step string) (int, bool) {
	switch string(step[len(step)-1]) {
	case "s":
		delta, err := strconv.Atoi(step[:len(step)-1])
		if err != nil {
			return 0, false
		}
		return delta, true
	case "m":
		delta, err := strconv.Atoi(step[:len(step)-1])
		if err != nil {
			return 0, false
		}
		return delta * 60, true
	case "h":
		delta, err := strconv.Atoi(step[:len(step)-1])
		if err != nil {
			return 0, false
		}
		return delta * 3600, true
	}
	return 0, false
}

func (f Field) generateTimestamp() (int, error) {
	if f.Value != nil {
		return (f.Value).(int), nil
	}

	var timestamp time.Time
	if f.Params.Now {
		timestamp = time.Now()
	}

	if f.Params.Incremental {
		var new_value int
		if value, ok := C.Get(f.uuid); ok {
			time_step, ok := getTimeStep(f.Params.Step.(string))
			if !ok {
				return 0, fmt.Errorf("Field %s doesn't have a valid step", f.Name)
			}
			var step int
			if f.Params.Sigma != 0 {
				step = int(distuv.Normal{Mu: float64(time_step), Sigma: f.Params.Sigma}.Rand())
			} else {
				step = time_step
			}
			new_value = value.(int) + step
		} else {
			if f.Params.Start == nil {
				new_value = int(time.Now().Unix())
			} else {
				new_value = int(f.Params.Start.(float64))
			}
		}
		C.Set(f.uuid, new_value, cache.NoExpiration)
		timestamp = time.Unix(int64(new_value), 0)
	}

	switch f.Params.Precision {
	case PARAMS_PRECISION_S:
		return int(timestamp.Unix()), nil
	case PARAMS_PRECISION_MS:
		return int(timestamp.UnixMilli()), nil
	case PARAMS_PRECISION_US:
		return int(timestamp.UnixMicro()), nil
	case PARAMS_PRECISION_NS:
		return int(timestamp.UnixNano()), nil
	default:
		return int(timestamp.Unix()), nil
	}

}

func (f Field) generateDatetime() (string, error) {
	if f.Value != nil {
		return (f.Value).(string), nil
	}

	var err error

	timestamp, err := f.generateTimestamp()
	if err != nil {
		return "", nil
	}

	datetime := time.Unix(int64(timestamp), 0)
	if f.Params.Format != "" {
		return datetime.Format(f.Params.Format), nil
	}

	return "", fmt.Errorf("Field %s doesn't specify a format", f.Name)
}

func (f Field) generateValue() (interface{}, error) {
	var value interface{}
	var err error

	switch f.Type {
	case UUID:
		value, err = f.generateUUID()
	case STRING:
		value, err = f.generateString()
	case INTEGER:
		value, err = f.generateInt()
	case FLOAT:
		value, err = f.generateFloat()
	case DATETIME:
		value, err = f.generateDatetime()
	case TIMESTAMP:
		value, err = f.generateTimestamp()
	case OBJECT:
		if f.Fields != nil {
			value, err = GenerateMessage(f.Fields)
		} else {
			value, err = nil, fmt.Errorf("Field %s doesn't include `FIELDS`", f.Name)
		}
	}
	return value, err
}

func GenerateMessage(fields *[]Field) (map[string]interface{}, error) {
	message := make(map[string]interface{})
	var err error

	for _, field := range *fields {

		if field.Nullable != 0 {
			if getRandomProbability() < field.Nullable {
				message[field.Name] = nil
				continue
			}
		}

		if field.Repeated {
			length := getRandomLength()
			objects := make([]interface{}, length)
			for i := 0; i < length; i++ {
				objects[i], err = field.generateValue()
			}
			message[field.Name], err = objects, err
			continue
		}

		message[field.Name], err = field.generateValue()
		if err != nil {
			panic(err)
		}

	}

	return message, err
}
