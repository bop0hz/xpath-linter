package lint

import (
	"github.com/antchfx/xmlquery"
	"strings"
	"testing"
)

const (
	tagExists = "<precondition></precondition><port>7777</port>"
)

func TestLookupIfTagExistsAndShouldbe(t *testing.T) {
	doc, _ := xmlquery.Parse(strings.NewReader(tagExists))
	rule := Rule{Contain: "/port"}
	found, node, err := rule.Lookup(doc)
	if !found {
		t.Fatal("Method does not find the node")
	}
	if node.InnerText() != "7777" {
		t.Fatal("Incorrect node text")
	}
	if err != nil {
		t.Fatalf("Method returns an error: %s", err)
	}
}

func TestLookupIfTagNotExistsAndShouldnot(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Contain: "/ip"}
	found, node, err := rule.Lookup(doc)
	if found {
		t.Fatal("Method finds something, but must not")
	}
	if node != nil {
		t.Fatal("Method returns something, but must not")
	}
	if err != nil {
		t.Fatalf("Method returns an error: %s", err)
	}
}

func TestCheckIfConditionExists(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Having: "/precondition"}
	suitable, err := rule.CheckIfCondition(doc)
	if !suitable {
		t.Fatal("Existing condition is not valid")
	}
	if err != nil {
		t.Fatalf("Method returns an error: %s", err)
	}
}

func TestCheckIfConditionNotExist(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Having: "/notprecondition"}
	suitable, err := rule.CheckIfCondition(doc)
	if suitable {
		t.Fatal("Not existing condition marks as suitable")
	}
}

func TestValidateShouldbeAndExist(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Contain: "/port", Must: true}
	valid, node, err := rule.Validate(doc)
	if !valid {
		t.Fatal("Case is not valid")
	}
	if node == nil {
		t.Fatal("Method returns nil node")
	}
	if err != nil {
		t.Fatalf("Method returns an error: %s", err)
	}
}

func TestValidateShouldBeAndNotExist(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Contain: "/ip", Must: true}
	valid, node, err := rule.Validate(doc)
	if valid {
		t.Fatal("Case is valid")
	}
	if node != nil {
		t.Fatal("Method returns not nil node")
	}
	if err == nil {
		t.Fatalf("Method does not return an error")
	}
}

func TestValidateShouldNotBeAndExist(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Contain: "/port", Must: false}
	valid, node, err := rule.Validate(doc)
	if valid {
		t.Fatal("Case is valid")
	}
	if node == nil {
		t.Fatal("Method returns nil node")
	}
	if err == nil {
		t.Fatal("Method does not return an error")
	}
}

func TestValidateShouldNotBeAndNotExist(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Contain: "/ip", Must: false}
	valid, node, err := rule.Validate(doc)
	if !valid {
		t.Fatal("Case is not valid")
	}
	if node != nil {
		t.Fatal("Method returns node")
	}
	if err != nil {
		t.Fatalf("Method returns an error: %s", err)
	}
}

func TestCheckShouldBe(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Having: "/precondition", Contain: "/port", Must: true}
	valid, node, err := rule.Check(doc)
	if !valid {
		t.Fatal("Case is not valid")
	}
	if node == nil {
		t.Fatal("Method returns nil node")
	}
	if err != nil {
		t.Fatalf("Method returns an error: %s", err)
	}
}

func TestCheckShouldNotBe(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Having: "/precondition", Contain: "/port", Must: false}
	valid, node, err := rule.Check(doc)
	if valid {
		t.Fatal("Case is valid")
	}
	if node == nil {
		t.Fatal("Method returns nil node")
	}
	if err == nil {
		t.Fatalf("Method does not return an error")
	}
}

func TestCheckWithIfConditionFalse(t *testing.T) {
	doc, err := xmlquery.Parse(strings.NewReader(tagExists))
	if err != nil {
		t.Fatal("Could not parse xml")
	}
	rule := Rule{Having: "/notprecondition", Contain: "/port", Must: true}
	valid, node, err := rule.Check(doc)
	if !valid {
		t.Fatal("Case is invalid")
	}
	if node != nil {
		t.Fatal("Method returns node")
	}
	if err != nil {
		t.Fatalf("Method returns an error")
	}
}
