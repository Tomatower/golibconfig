This is an golang implementation for [libconfig](https://github.com/hyperrealm/libconfig).

It utilizes the native libconfig implementation, so this has to be installed on
your system.

# Exemplary usage

The API orients itself along libconfig++.

Consider the following config file:

```config
teststring = "test";
testint = 1;
testfloat = 0.1;
testbool = false;

testgroup : {
	testingroup = "in group";
}
```

this can be interfaced like the following

```go

config := libconf.NewConfig()
err := config.ReadFile("configfile.conf")
if err != nil {
	panic(err)
}

intval, err := config.LookupInt("testint")

group := config.Lookup("testgroup")
teststring, err := group.LookupString("testinggroup")
```
