package main

import (
	"context"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"os/user"
	"strings"

	"github.com/digitalocean/godo"
	"golang.org/x/oauth2"
)

type TokenSource struct {
	AccessToken string
}

func (t *TokenSource) Token() (*oauth2.Token, error) {
	token := &oauth2.Token{
		AccessToken: t.AccessToken,
	}
	return token, nil
}

func main() {
	pat := flag.String("token", "NONE", "Digital Ocean access token")
	clusterID := flag.String("cluster", "NONE", "Digital Ocean cluster ID")
	flag.Parse()
	usr, err := user.Current()

	tokenSource := &TokenSource{
		AccessToken: *pat,
	}

	oauthClient := oauth2.NewClient(oauth2.NoContext, tokenSource)
	client := godo.NewClient(oauthClient)

	ctx := context.TODO()

	config, _, err := client.Kubernetes.GetKubeConfig(ctx, *clusterID)
	if err != nil {
		log.Fatal(err)
	}

	kubeConfigFile := string(config.KubeconfigYAML)
	var str strings.Builder
	str.WriteString(usr.HomeDir)
	str.WriteString("/.kube")
	_ = os.Mkdir(str.String(), 0770)
	str.WriteString("/config")
	err = ioutil.WriteFile(str.String(), []byte(kubeConfigFile), 0640)
	if err != nil {
		log.Fatal(err)
	}
}
