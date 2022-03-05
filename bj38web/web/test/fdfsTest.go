package main
import "github.com/tedcy/fdfs_client"
func main () {
	fdfs_client.NewClientWithConfig("/etc/fdfs/client.conf")
}
