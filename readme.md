# Clails - cloud made easy

Clails is an opinionated toolkit which takes simple project definition in YML file and generates complete cloud
 templates (like [AWS CloudFormation](https://aws.amazon.com/cloudformation)) from it. 
 Think "*[Rails](https://rubyonrails.org) but for cloud*".
 
## Usage

In order provision new project, create Clails project file and save it as `clails.yml`: 

```
name: MyProject
services:
  - type: kafka
```

Then execute the following command in the same directory:

```
$ clails deploy
```
 
 ## License
 
 This project is distributed under [Apache 2.0 license](http://www.apache.org/licenses/LICENSE-2.0.html).