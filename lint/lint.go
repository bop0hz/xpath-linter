package lint

import (
	"fmt"

	"github.com/antchfx/xmlquery"
	"github.com/bop0hz/xml-linter/report"
)

// Rule object to check
type Rule struct {
	Name    string `yaml: "name"`
	Targets string `yaml: "targets"`
	Having  string `yaml: "having"`
	Contain string `yaml: "contain"`
	Must    bool   `yaml: "must"`
}

// CheckIfCondition checks if condition in rule for filtering document as suitable or not
func (r *Rule) CheckIfCondition(doc *xmlquery.Node) (suitable bool, err error) {
	if r.Having != "" {
		node, err := xmlquery.Query(doc, r.Having)
		if err != nil {
			err = fmt.Errorf("Could not parse query expression: %s", err)
			return false, err
		}
		if node == nil {
			return false, nil
		}
	}
	return true, nil
}

// Lookup for xml node
func (r *Rule) Lookup(doc *xmlquery.Node) (found bool, node *xmlquery.Node, err error) {
	node, err = xmlquery.Query(doc, r.Contain)
	if err != nil {
		err = fmt.Errorf("Could not parse query expression: %s", err)
		return false, nil, err
	}
	if node != nil {
		return true, node, nil
	}
	return false, nil, nil
}

// Validate Shouldbe parameter
func (r *Rule) Validate(doc *xmlquery.Node) (valid bool, node *xmlquery.Node, err error) {
	found, node, err := r.Lookup(doc)
	if r.Must && !found {
		err = fmt.Errorf("Expected element not found")
		return false, nil, err
	}
	if !r.Must && found {
		err = fmt.Errorf("Unexpected element found")
		return false, node, err
	}
	return true, node, nil
}

// Check whole linting rule
func (r *Rule) Check(doc *xmlquery.Node) (valid bool, node *xmlquery.Node, err error) {
	suitable, err := r.CheckIfCondition(doc)
	if err != nil {
		fmt.Errorf("Error occured during precondition check: %s, %v", r.Having, err)
		return false, nil, err
	}
	if suitable {
		valid, node, err := r.Validate(doc)
		if err != nil {
			return false, node, err
		}
		if valid {
			return true, node, nil
		}
	}
	return true, node, nil
}

// CheckAll checks all rules in xml document
func CheckAll(doc *xmlquery.Node, rules []*Rule, r report.Reporter) (results []string, err error) {
	for _, rule := range rules {
		targets := "/"
		if rule.Targets != "" {
			targets = rule.Targets
		}
		nodes, err := xmlquery.QueryAll(doc, targets)
		if err != nil {
			fmt.Errorf("Could not query nodes: %s", err)
			return nil, err
		}
		for i, node := range nodes {
			valid, validatedNode, err := rule.Check(node)
			if !valid && validatedNode == nil {
				results = append(results, r.Compile(rule.Name, i, err.Error(), rule.Contain))
			}
			if !valid && validatedNode != nil {
				results = append(results, r.Compile(rule.Name, i, err.Error(), validatedNode.OutputXML(true)))
			}
		}
	}
	return results, nil
}
