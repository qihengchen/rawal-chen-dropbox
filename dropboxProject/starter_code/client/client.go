package main

import (
	"fmt"
	"os"

	"../internal"
	"../lib/support/client"
	"../lib/support/rpc"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Fprintf(os.Stderr, "Usage: %v <server>\n", os.Args[0])
		os.Exit(1)
	}

	// EXAMPLE CODE
	//
	// This code is meant as an example of how to use
	// our framework, not as stencil code. It is not
	// meant as a suggestion of how you should write
	// your application.

	server := rpc.NewServerRemote(os.Args[1])

	// Examples of calling various functions on the server
	// over RPC.

	var retInt int
	err := server.Call("add", &retInt, 2, 4)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calling method add: %v\n", err)
		return
	} else {
		fmt.Printf("add(2, 4): %v\n", retInt)
	}

	err = server.Call("mult", &retInt, 2, 4)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calling method mult: %v\n", err)
		return
	}
	fmt.Printf("mult(2, 4): %v\n", retInt)

	err = server.Call("noOp", nil)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error calling method noOp: %v\n", err)
		return
	}
	fmt.Println("noOp()")

	// An example of how you might run a basic client.

	// In a real client, you'd have to first authenticate the user
	// (note that we don't provide any support code for this,
	// including the command-line interface). Once you the user
	// is authenticated, the client object (of the Client type
	// in this example, but it can be anything that implements
	// the client.Client interface) should somehow keep hold of
	// session information so that future requests (initiated
	// by methods being called on the object) can be authenticated.

	c := Client{server}
	err = client.RunCLI(&c)
	if err != nil {
		// don't actually log the error; it's already been
		// printed by client.RunCLI
		os.Exit(1)
	}
}

// An implementation of a basic client to match the example server
// implementation. This client/server implementation is absurdly
// insecure, and is only meant as an example of how to implement
// the client.Client interface; it should not be taken as a suggestion
// of how to design your client or server.
type Client struct {
	server *rpc.ServerRemote
}

func (c *Client) Upload(path string, body []byte) (err error) {
	var ret string
	err = c.server.Call("upload", &ret, path, body)
	if err != nil {
		return client.MakeFatalError(err)
	}
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (c *Client) Download(path string) (body []byte, err error) {
	var ret internal.DownloadReturn
	err = c.server.Call("download", &ret, path)
	if err != nil {
		return nil, client.MakeFatalError(err)
	}
	if ret.Err != "" {
		return nil, fmt.Errorf(ret.Err)
	}
	return ret.Body, nil
}

func (c *Client) List(path string) (entries []client.DirEnt, err error) {
	var ret internal.ListReturn
	err = c.server.Call("list", &ret, path)
	if err != nil {
		return nil, client.MakeFatalError(err)
	}
	if ret.Err != "" {
		return nil, fmt.Errorf(ret.Err)
	}
	var ents []client.DirEnt
	for _, e := range ret.Entries {
		ents = append(ents, e)
	}
	return ents, nil
}

func (c *Client) Mkdir(path string) (err error) {
	var ret string
	err = c.server.Call("mkdir", &ret, path)
	if err != nil {
		return client.MakeFatalError(err)
	}
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (c *Client) Remove(path string) (err error) {
	var ret string
	err = c.server.Call("remove", &ret, path)
	if err != nil {
		return client.MakeFatalError(err)
	}
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}

func (c *Client) PWD() (path string, err error) {
	var ret internal.PWDReturn
	err = c.server.Call("pwd", &ret)
	if err != nil {
		return "", client.MakeFatalError(err)
	}
	if ret.Err != "" {
		return "", fmt.Errorf(ret.Err)
	}
	return ret.Path, nil
}

func (c *Client) CD(path string) (err error) {
	var ret string
	err = c.server.Call("cd", &ret, path)
	if err != nil {
		return client.MakeFatalError(err)
	}
	if ret != "" {
		return fmt.Errorf(ret)
	}
	return nil
}
