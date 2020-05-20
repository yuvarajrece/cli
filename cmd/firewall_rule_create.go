package cmd

import (
	"fmt"
	"strings"

	"github.com/civo/civogo"
	"github.com/civo/cli/config"
	"github.com/civo/cli/utility"

	"os"

	"github.com/spf13/cobra"
)

var protocol, startPort, endPort, direction, label string
var cidr []string

var firewallRuleCreateCmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"new", "add"},
	Short:   "Create a new firewall rule",
	Args:    cobra.MinimumNArgs(1),
	Example: "civo firewall rule create FIREWALL_NAME/FIREWALL_ID [flags]",
	Run: func(cmd *cobra.Command, args []string) {
		client, err := config.CivoAPIClient()
		if err != nil {
			utility.Error("Creating the connection to Civo's API failed with %s", err)
			os.Exit(1)
		}

		firewall, err := client.FindFirewall(args[0])
		if err != nil {
			utility.Error("Unable to find the firewall for your search %s", err)
			os.Exit(1)
		}

		newRuleConfig := &civogo.FirewallRuleConfig{
			FirewallID: firewall.ID,
			Protocol:   protocol,
			StartPort:  startPort,
			Cidr:       cidr,
			Direction:  direction,
			Label:      label,
		}

		if endPort == "" {
			newRuleConfig.EndPort = startPort
		} else {
			newRuleConfig.EndPort = endPort
		}

		rule, err := client.NewFirewallRule(newRuleConfig)
		if err != nil {
			utility.Error("Unable to create the new firewall rule %s", err)
			os.Exit(1)
		}

		ow := utility.NewOutputWriterWithMap(map[string]string{"ID": rule.ID, "Name": rule.Label})

		switch outputFormat {
		case "json":
			ow.WriteSingleObjectJSON()
		case "custom":
			ow.WriteCustomOutput(outputFields)
		default:
			if newRuleConfig.EndPort == newRuleConfig.StartPort {
				fmt.Printf("Created a firewall rule called %s allowing access to port %s from %s with ID %s\n", utility.Green(rule.Label), utility.Green(newRuleConfig.StartPort), utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
			} else {
				fmt.Printf("Created a firewall rule called %s allowing access to ports %s-%s from %s with ID %s\n", utility.Green(rule.Label), utility.Green(newRuleConfig.StartPort), utility.Green(newRuleConfig.EndPort), utility.Green(strings.Join(newRuleConfig.Cidr, ", ")), rule.ID)
			}
		}
	},
}
