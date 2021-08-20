# Serenity

A sandox environment for testing Prometheus and Alertmanager
inhibition rules.


## Overview
There are three major components: Serenity, Prometheus and Alertmanager.

Serenity is a simple HTTP service that exposes two gauges: `foo` and `bar`.
Both gauges have initial value of 0, but we can flip them to 1
by sending HTTP requests like `GET /foo/on` and `GET /bar/on` respectively.

Serenity exposes its metrics via `GET /metrics`
and they are getting scraped by Prometheus.

Prometheus is configured with two alerts: `FooIsOnFire` and `BarIsOnFire`;
they fire as soon as the corresponding metric goes above zero.
Those alerts are sent to Alertmanager.

Alertmanager is configured with an inhibition rule that prevents `BarIsOnFire` from firing
if `FooIsOnFire` is already firing, but only if their "cluster" label matches.

As this is a sandbox, Alertmanager is not sending any alerts downstream,
but we can still see which alerts are currently firing by inspecting
Alertmanager's state using its Web UI.


## Requirements

* [Docker Desktop (Mac or Windows)](https://www.docker.com/products/docker-desktop)
* [skaffold](https://skaffold.dev/)


## One-Time Setup

1. Enable Kubernetes in your Docker Desktop preferences.
1. **Triple-check that your current kubectl context is called `docker-desktop`**:
	```
	$ kubectl config current-context
	docker-desktop
	```

	**This is extremely important as there currently are no safeguards for accidentally breaking production if your kubectl context is misconfigured!**
1. Apply "setup" manifests:
    ```
    kubectl apply -f k8s/setup/
    ```

    This will install Prometheus Operator and create instances of Prometheus and Alertmanager for us to configure.

## Usage

1. Run `skaffold dev`.
    This will start Serenity and apply Prometheus and Alertmanager configuration.

1. Open [Prometheus Web UI](http://127.0.0.1:9090/) and
    [check out our two metrics](http://127.0.0.1:9090/graph?g0.expr=%7B__name__%3D~%22foo%7Cbar%22%7D&g0.tab=0&g0.stacked=0&g0.show_exemplars=0&g0.range_input=1h).

1. Send the following requests to Serenity and observe the change in metrics:
    ```
    $ curl http://127.0.0.1:8080/foo/on
    $ curl http://127.0.0.1:8080/bar/on
    ```

    FYI you can also use `GET /foo/off` and `GET /bar/off` to reset the gauges to zero.

1. [Check the list of alerts that are firing _from Prometheus' perspective_](http://127.0.0.1:9090/alerts)

    You should see both alerts firing.

1. Open [Alertmanaget Web UI](http://127.0.0.1:9093/#/alerts) and check the list of alerts
    that are firing _from Alertmanager's perspective_.

    You should see that only `FooIsOnFire` is firing.

    If you click on the "Inhibited" checkbox in the upper right corner,
    you should be able to see both alerts, but `BarIsOnFire` will be marked as inhibited.

1. (Bonus step) You can check out how label matching works on inhibition rules by
    opening `k8s/service-monitor.yaml` and changing the line that says
    ```
          replacement: eu-west-1
    ```
    to say:
    ```
          replacement: eu-west-2
    ```

    Once you save the file, the change will be automatically applied by `skaffold`
    and shortly you should see both alerts firing in Alertmanager as they are now
    coming from different "regions".

## Cleanup

```
$ kubectl delete -f k8s/setup/
```

There will be a couple of "NotFound" errors which are safe to ignore.
