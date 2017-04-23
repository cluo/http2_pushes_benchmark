package main

import (
	"fmt"
	"os"
)

const jsFilesNumber = 1000

//Generates a html file and javascripts for serv to send
func generateResources() {
	err := os.Mkdir("res", 0777)
	if err != nil {
		panic(err)
	}
	html := "<html>\n	<head>\n"

	//Create js files and add a <script> tag for each of them
	for i := 0; i <= jsFilesNumber; i++ {
		filename := fmt.Sprintf("f%03d.js", i)
		os.Remove("res/" + filename)
		f, err := os.Create("res/" + filename)
		if err != nil {
			panic(err)
			return
		}
		f.Write([]byte(fmt.Sprintf("var a;")))
		f.Close()
		html += fmt.Sprintf("		<script src='res/%s' ></script>\n", filename)
	}
	html += `	</head>
	<body>
		<script type="text/javascript">
           window.onload = function () {
    			var loadTime = window.performance.timing.domContentLoadedEventEnd-window.performance.timing.navigationStart; 
    			console.log('Page load time is '+ loadTime + 'ms');
		   }
		</script>
	</body>
	</html>`
	os.Remove("res/index.html")
	f, err := os.Create("res/index.html")
	if err != nil {
		panic(err)
	}
	f.Write([]byte(html))
	f.Close()
}
