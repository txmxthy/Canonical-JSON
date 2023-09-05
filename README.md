
<h1 align="center">Welcome to Canon-JSON 👋</h1>
<p align="center">
 <a href="#license">
	<img src="https://img.shields.io/badge/License-MIT-blue?style=for-the-badge" alt="License"></a>
<a href="https://github.com/txmxthy/canon-json/issues">
	<img src="https://img.shields.io/github/issues/txmxthy/canon-json?style=for-the-badge" alt="issues - canon-json"></a>
<a href="https://github.com/txmxthy/canon-json">
	<img src="https://img.shields.io/github/stars/txmxthy/canon-json?style=for-the-badge" alt="stars - canon-json"></a>
<a href="https://github.com/txmxthy/canon-json">
	<img src="https://img.shields.io/github/forks/txmxthy/canon-json?style=for-the-badge" alt="forks - canon-json"></a>
</p>




<p align="center">
	<a href="https://golang.org">
		<img src="https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white" alt="Go"></a>
</p>


## About
This project is home to my 3 hour 3am json canonical spec implementation pieced together from various resources (listed below)


### Key Features
- Implements Json Canonical Spec (more or less) according to http://gibson042.github.io/canonicaljson-spec/


## Usage

To test the project, just download the code and run the main_test.go file with

```bash
  go run main_test.go
```
If you want to test breaking it I reccomend trying to flip the !first boolean in canonicalize then running the tests again. 


## Todo
Change tests to use data file instead of hardcoded cases.
  
## Contributing

Contributions are always welcome!


## References I used :)
http://gibson042.github.io/canonicaljson-spec/

https://github.com/gibson042/canonicaljson-go

https://github.com/goccy/go-json

  
