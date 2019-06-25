package testrepos

type SSHRepo struct {
	URL, SSHPrivateKey    string
	InsecureIgnoreHostKey bool
}

// why isn't this the same repo as the HTTPS one? because you could accidentally use the username/password to access it
// and that would invalidate all test for SSH
var SSHTestRepo = SSHRepo{
	URL: "git@gitlab.com:alex_collins/private-repo.git",
	SSHPrivateKey: `-----BEGIN OPENSSH PRIVATE KEY-----
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
-----END OPENSSH PRIVATE KEY-----
`,
	InsecureIgnoreHostKey: true,
}
