package gode

import "os/exec"

// RunScript runs a given script in node
// Returns an *os/exec.Cmd instance
func (c *Client) RunScript(script string) *exec.Cmd {
	cmd := exec.Command(c.NodePath, "-e", script)
	cmd.Env = []string{"NODE_PATH=" + c.ModulesPath}
	return cmd
}
