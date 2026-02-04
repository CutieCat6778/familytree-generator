package generator

type LifeExpectancyMode string

const (
	LifeExpectancyTotal    LifeExpectancyMode = "total"
	LifeExpectancyFemale   LifeExpectancyMode = "female"
	LifeExpectancyMale     LifeExpectancyMode = "male"
	LifeExpectancyByGender LifeExpectancyMode = "by_gender"
)

func ParseLifeExpectancyMode(value string) LifeExpectancyMode {
	switch LifeExpectancyMode(value) {
	case LifeExpectancyFemale, LifeExpectancyMale, LifeExpectancyByGender, LifeExpectancyTotal:
		return LifeExpectancyMode(value)
	default:
		return LifeExpectancyTotal
	}
}
