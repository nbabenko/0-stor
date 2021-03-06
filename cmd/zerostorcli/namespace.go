package main

import (
	"fmt"

	"github.com/urfave/cli"
	"github.com/zero-os/0-stor/client/itsyouonline"
)

func createNamespace(c *cli.Context) error {
	iyoCl, err := getNamespaceManager(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if len(c.Args()) < 1 {
		return cli.NewExitError(fmt.Errorf("need to give the name of the namespace to create"), 1)
	}

	namespace := c.Args().First()
	if err := iyoCl.CreateNamespace(namespace); err != nil {
		return cli.NewExitError(fmt.Errorf("creation of namespace %s failed: %v", namespace, err), 1)
	}
	fmt.Printf("Namespace %s created\n", namespace)

	return nil
}

func deleteNamespace(c *cli.Context) error {
	iyoCl, err := getNamespaceManager(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	if len(c.Args()) < 1 {
		return cli.NewExitError(fmt.Errorf("need to give the name of the namespace to create"), 1)
	}

	namespace := c.Args().First()
	if err := iyoCl.DeleteNamespace(namespace); err != nil {
		return cli.NewExitError(err, 1)
	}
	fmt.Printf("Namespace %s deleted\n", namespace)
	return nil
}

func setACL(c *cli.Context) error {
	iyoCl, err := getNamespaceManager(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	namespace := c.String("namespace")
	user := c.String("user")
	currentPermission, err := iyoCl.GetPermission(namespace, user)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("fail to retrieve permission : %v", err), 1)
	}

	requestedPermission := itsyouonline.Permission{
		Read:   c.Bool("r"),
		Write:  c.Bool("w"),
		Delete: c.Bool("d"),
		Admin:  c.Bool("a"),
	}

	// remove permission if needed
	toRemove := itsyouonline.Permission{
		Read:   currentPermission.Read && !requestedPermission.Read,
		Write:  currentPermission.Write && !requestedPermission.Write,
		Delete: currentPermission.Delete && !requestedPermission.Delete,
		Admin:  currentPermission.Admin && !requestedPermission.Admin,
	}
	if err := iyoCl.RemovePermission(namespace, user, toRemove); err != nil {
		return cli.NewExitError(err, 1)
	}

	// Give requested permission
	if err := iyoCl.GivePermission(namespace, user, requestedPermission); err != nil {
		return cli.NewExitError(err, 1)
	}

	return nil
}

func getACL(c *cli.Context) error {
	iyoCl, err := getNamespaceManager(c)
	if err != nil {
		return cli.NewExitError(err, 1)
	}

	namespace := c.String("namespace")
	user := c.String("user")
	perm, err := iyoCl.GetPermission(namespace, user)
	if err != nil {
		return cli.NewExitError(fmt.Errorf("fail to retrieve permission : %v", err), 1)
	}
	fmt.Printf("User %s:\n", user)
	fmt.Printf("Read: %v\n", perm.Read)
	fmt.Printf("Write: %v\n", perm.Write)
	fmt.Printf("Delete: %v\n", perm.Delete)
	fmt.Printf("Admin: %v\n", perm.Admin)

	return nil
}
