# XML-linter

Simple, yet flexible XPath powered XML linter.

This tool implies you are familiar with [XPath](https://en.wikipedia.org/wiki/XPath) (XML Path Language)

## Modes
* CI mode designed for Continuous Integration or batch checks - evaluates rules from a configuration file
* Ad-hoc query - evaluates rules passed by command line flags

## Getting started
In CI mode the tool should be used in collaboration with exec option of find utility. For example, if we are going to check all xml in Regression directory against rules in config.yaml:

```bash
find ./Regression/ -type f -name "*.xml" -exec ./xpath-linter -ci -cfg ./config.yaml {} \;
```

In ad-hoc mode, e.g. we are looking for services which do not have generic variables:
```bash
find ./Regression/ -type f -name "*.xml" -exec ./xpath-linter -contain //variables {} \;
```


## Linting rules example

config.yaml:
```yaml
---
rules:
  - name: No variables in settings
    targets: /settings
    must: yes
    contain: //variables

  - name: Username in settings having type=Client is not empty
    having: //type[text()="Client"]
    must: no
    contain: //username[text()]

  - name: No nodes with value > 1
    targets: //node
    must: no
    contain: //node[text()>1]

  - name: Empty tag
    targets: //node
    must: no
    contain: //node[not(text())]

  - name: Tag contains bad word
    must: no
    contain: //tag[contains(text(), "bad word")]
```


## Command line interface

```shell
Usage: ./release/xpath-linter [FLAGS] [FILE]
  -cfg string
    	File with linting rules (default "config.yaml")
  -ci
    	Continous integration mode. Evalutates rules from the config against the target file
  -contain string
    	XPath query to evaluate. Cheatsheet https://devhints.io/xpath (default "/")
  -having string
    	Optional XPath condition which must be true to proceed with a query. Works like a filter
  -must
    	Whether target node must have or not 'contain' nodes (default true)
  -targets string
    	XPath query to find target nodes, default is root node (default "/")
  -version
    	Show app version
```

## XPath cheatsheet
[https://devhints.io/xpath#operators](https://devhints.io/xpath#operators)

### Build
`make`
