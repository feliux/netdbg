## Under Construction

<img src="images/netdbg.png"
     alt="Gopher netdbg"
	style="float:left; margin-right:10px; width:200px; height:200px;"
/>

Netdbg is a powerfull command line tool for debugging connectivity issues in enviroments where you do not have tools like ping, curl, netcat, dig, nslookp, etc, etc. So netdbg is the swiss knife that brings all of them together in one place.

### Why?

Think in a distroless container deployed in Kubernetes. How do you test connectivity issues from that container if the image just have the necessary binaries for be able to run your application. So... Maybe your application can no connect to other pod... Maybe you do not know about the subnet where the cluster is deployed... Maybe there are firewalls blocking your application... Maybe you do not set the correct proxy endpoint... etc etc. With netdbg you can execute ping, netcat and that all stuff necessary to check out why your application is failing.
