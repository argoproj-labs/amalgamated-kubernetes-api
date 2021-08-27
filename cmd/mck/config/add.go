package config

import (
	context "context"
	"fmt"
	"github.com/argoproj-labs/multi-cluster-kubernetes/api/config"
	"github.com/spf13/cobra"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	clientcmdapi "k8s.io/client-go/tools/clientcmd/api"
	"k8s.io/client-go/util/homedir"
	"path/filepath"
)

func NewAddCommand() *cobra.Command {
	var (
		kubeconfig string
		namespace  string
	)
	cmd := &cobra.Command{
		Use: "add [CONTEXT_NAME]",
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()

			startingConfig, err := clientcmd.NewDefaultPathOptions().GetStartingConfig()
			if err != nil {
				return err
			}

			if len(args) == 1 {
				startingConfig.CurrentContext = args[0]
			}

			clientConfig := clientcmd.NewNonInteractiveDeferredLoadingClientConfig(&clientcmd.ClientConfigLoadingRules{ExplicitPath: kubeconfig}, &clientcmd.ConfigOverrides{})
			restConfig, err := clientConfig.ClientConfig()
			if err != nil {
				return err
			}
			if namespace == "" {
				namespace, _, err = clientConfig.Namespace()
				if err != nil {
					return err
				}
			}

			secretsInterface := kubernetes.NewForConfigOrDie(restConfig).CoreV1().Secrets(namespace)

			err = clientcmdapi.MinifyConfig(startingConfig)
			if err != nil {
				return err
			}
			err = config.New(secretsInterface).Add(ctx, startingConfig)
			if err != nil {
				return err
			}

			fmt.Printf("context %q added\n", startingConfig.CurrentContext)

			return nil
		},
	}
	cmd.Flags().StringVar(&kubeconfig, "kubeconfig", filepath.Join(homedir.HomeDir(), ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	cmd.Flags().StringVarP(&namespace, "namespace", "n", "", "namespace")
	return cmd
}
