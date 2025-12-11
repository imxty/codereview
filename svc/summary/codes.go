package summary

import "fmt"

//go:generate ./summary_gen.sh
//go:generate ./risk_gen.sh
//go:generate ./physique_gen.sh

// 中医报告小结
type Summary struct {
	Label   string `yaml:"label"`
	Summary string `yaml:"summary"`
}

func (c Summary) CommentLine() string {
	return fmt.Sprintf("// %s", c.Label)
}

// 疾病风险
type Risk struct {
	Label string `yaml:"label"`
}

func (c Risk) CommentLine() string {
	return fmt.Sprintf("// %s", c.Label)
}

// 体质
type Physique struct {
	Label string `yaml:"label"`
}

func (c Physique) CommentLine() string {
	return fmt.Sprintf("// %s", c.Label)
}
