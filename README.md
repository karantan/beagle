# Beagle
Beagle is a lightweight tool for finding long-running processes and potentially isolating
them. It also reports its findings on Slack (if configured).

The beagle is a breed of small scent hound. It was developed primarily for hunting rabbits.
They possess a great sense of smell and superior tracking instincts, and this is precisely
what we need when we want to find suspicious processes.

PHP-FPM child processes are usually short-lived, so if we see a PHP-FPM pool running
for several hours then this could indicate that there is some malicious code running
without us knowing about it.

Beagle (tool) will find it, report it and potentially isolate it in a separate cgroup
where it can be controlled until we figure out what to do with it.
