package mypackages

import (
	"bufio"
	"bytes"
	"html/template"
	"os"

)

// type product struct {
// 	Img         string
// 	Name        string
// 	Price       string
// 	Stars       float64
// 	Reviews     int
// 	Description string
// 	Values      string
// }

type product struct {
	JobName       string
	JobBackuptime string
	JBinc         string
	JobReptime    string
	JRinc         string
	JobArctime    string
	JAinc         string
}

func subtr(a, b float64) float64 {
	return a - b
}

func list(e ...float64) []float64 {
	return e
}

var data []product

func Addmetrics(name string, btime string, binc string, rtime string, rinc string, atime string, ainc string) {
	var temp product
	temp.JobName = name
	temp.JobBackuptime = btime
	temp.JBinc = binc
	temp.JobReptime = rtime
	temp.JRinc = rinc
	temp.JobArctime = atime
	temp.JAinc = ainc

	data = append(data, temp)

	//	fmt.Println(datanew)

}

func generator() {

	// data := []product{
	// 	{"images/1.png", "abcc", "$11.00", 4.0, 251, "Lorem ipsum dolor sit amet, consectetur adipiscing elit.", "abcddccc"},
	// 	// {"images/2.png", "onions", "$2.80", 5.0, 123, "Morbi sit amet erat vitae purus consequat vehicula nec sit amet purus."},
	// 	// {"images/3.png", "tomatoes", "$3.10", 4.5, 235, "Curabitur tristique odio et nibh auctor, ut sollicitudin justo condimentum."},
	// 	// {"images/4.png", "courgette", "$1.20", 4.0, 251, "Phasellus at leo a purus consequat ornare ac aliquam quam."},
	// 	// {"images/5.png", "broccoli", "$3.80", 3.5, 123, "Maecenas sed ante sagittis, dapibus dui quis, hendrerit orci."},
	// 	// {"images/6.png", "potatoes", "$3.00", 2.5, 235, "Vivamus malesuada est et tellus porta, vel consectetur orci dapibus."},
	// }
	//	addmetrics("abc", "b", 1, "r", 2, "a", 3)

	allFiles := []string{"content.tmpl", "footer.tmpl", "header.tmpl", "page.tmpl"}

	var allPaths []string
	for _, tmpl := range allFiles {
		allPaths = append(allPaths, "./templates/"+tmpl)
	}

	templates := template.Must(template.New("").Funcs(template.FuncMap{"subtr": subtr, "list": list}).ParseFiles(allPaths...))

	var processed bytes.Buffer
	templates.ExecuteTemplate(&processed, "page", data)

	outputPath := "./static/index.html"
	f, _ := os.Create(outputPath)
	w := bufio.NewWriter(f)
	w.WriteString(string(processed.Bytes()))
	w.Flush()

}
