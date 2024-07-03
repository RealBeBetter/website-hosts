package core

const startTag = "# GitHub Start\r\n"
const endTag = "# GitHub End\r\n"

type HostChan struct {
	Domain string
	Ip     string
	Err    error
}

var domains = []string{
	"github.com",
	//"gist.github.com",
	//"assets-cdn.github.com",
	//"raw.githubusercontent.com",
	//"gist.githubusercontent.com",
	//"cloud.githubusercontent.com",
	//"camo.githubusercontent.com",
	//"avatars.githubusercontent.com",
	//"avatars0.githubusercontent.com",
	//"avatars1.githubusercontent.com",
	//"avatars2.githubusercontent.com",
	//"avatars3.githubusercontent.com",
	//"avatars4.githubusercontent.com",
	//"avatars5.githubusercontent.com",
	//"avatars6.githubusercontent.com",
	//"avatars7.githubusercontent.com",
	//"avatars8.githubusercontent.com",
	//"github.githubassets.com",
}
