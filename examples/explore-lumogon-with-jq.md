# Querying Lumogon with JQ

Lumogon is a low level tool that exposes JSON output when run. But that doesn't
mean you can't ask specific questions of the data returned. Lets see some examples
of doing just that with the excellent [jq](https://stedolan.github.io/jq/) JSON
processor.

## A handy alias

Packaging up Lumogon as a Docker container has a number of advantages,
but the verbose command to run it can get old quickly.

We recommend creating an alias in your shell of choice to skip the long command
line invocation.

```
alias lumogon="docker run --rm -v /var/run/docker.sock:/var/run/docker.sock puppet/lumogon"
```

This alias allows you to run all of the examples in this post with
`lumogon` rather than the full `docker` incantation. We’ll use this alias for
the following examples.

## Looking at a containers labels

Let’s look at a simple example to begin with. Let’s view the labels for
a specific container:

```
$ lumogon scan fixtures_ubuntu-xenial_1 | jq ".containers[].capabilities.label.payload"
"payload": {
  "com.docker.compose.config-hash":
"600c3e117e3a2dfeadac8bec2680b40f71dcc2fe8dae8d402432131df2d59646",
  "com.docker.compose.container-number": "1",
  "com.docker.compose.oneoff": "False",
  "com.docker.compose.project": "fixtures",
  "com.docker.compose.service": "ubuntu-xenial",
  "com.docker.compose.version": "1.11.2"
}
```

That’s not hugely interesting, you could do the same with `docker
inspect -f "{{json .Config.Labels }}"`. So let’s do something more
interesting that’s not available from the native Docker tools.

## Listing packages installed in a container

Many containers are using Linux distributions as the user space for their
containers, which means they contain a list of packages at specific
versions. Lumogon collects that information (from the dpkg, rpm, and apk package
managers) and makes it available to us. For instance:

```
$ lumogon scan fixtures_ubuntu-xenial_1 | jq ".containers[].capabilities.dpkg.payload" | head
{
  "adduser": "3.113+nmu3ubuntu4",
  "apt": "1.2.18",
  "base-files": "9.4ubuntu4.3",
  "base-passwd": "3.5.39",
  "bash": "4.3-14ubuntu1.1",
  "bsdutils": "1:2.27.1-6ubuntu3.2",
  "coreutils": "8.25-2ubuntu2",
  "dash": "0.5.8-2.1ubuntu2",
  "debconf": "1.5.58ubuntu1",
```

Here we start to see some of the value of Lumogon, it’s a single tool
which can query information about your container from the outside (like
labels which come from the Docker API) and the inside (like installed
packages or the user space OS details), without knowing anything about
the container in question beforehand.

One more example before we move on, let’s grab the version of bash
installed in a specific container:

```
$ lumogon scan fixtures_ubuntu-xenial_1 | jq ".containers[].capabilities.dpkg.payload.bash"
"4.3-14ubuntu1.1"
```

## Information about more than one container

All the examples so far have shown us gathering data about individual
containers. Let’s expand that and collect information from all of the
containers running on this host. To make this more interesting we’ll
then filter that information, producing a list showing the container
names and the version of bash installed via either rpm or dpkg.

```
$ lumogon scan | jq -r  '.containers[] | .container_name + "   " + .capabilities.dpkg.payload.bash + "    " + .capabilities.rpm.payload.bash'
/fixtures_debian-jessie_1   4.3-11+deb8u1
/fixtures_alpine_1
/fixtures_centos7_1       4.2.46-21.el7_3-x86_64
/fixtures_fedora_1       4.3.43-4.fc25-x86_64
/fixtures_debian-wheezy_1   4.2+dfsg-0.1+deb7u4
/fixtures_ubuntu-xenial_1   4.3-14ubuntu1.1
```

Another example, this time we will search through all of our containers for
those using a debian derivative.

```
$ lumogon scan | jq -r '.containers[] | select(.capabilities.host.payload.platformfamily == "debian") | .container_name'
/fixtures_debian-jessie_1
/fixtures_ubuntu-trusty_1
/fixtures_debian-wheezy_1
/fixtures_ubuntu-xenial_1
```

This type of inspection is much more powerful than simply relying on the parent
image for this information. We’re inspecting the file system itself
rather than relying on usefully named images.

Consider how you would do verify this information without Lumogon? It’s a
multi-step process that requires intimate implementation knowledge. With Lumogon
we want to set that information free so we can build tools to quickly answer
questions and solve real user problems.

If you're interested in these examples or have any questions about
Lumogon then head over to our [Slack channel](https://puppetcommunity.slack.com/messages/C5CT7GMKQ).
