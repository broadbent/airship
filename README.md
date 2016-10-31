# Airship

*A Combinatorial Auctioneer for Fog Resource Provisioning*

## Introduction

Airship is a tool designed to provide market-based resource allocation in Fog environments. This is intended to address the challenge of satisfying the requirements of multiple potential competing infrastructure users.

To do this, Airship uses a PAUSE combinatorial auction system. This is designed to be computationally tractable, a common pitfall of alternative techniques. More details can be found in [Chapter 6](https://doi.org/10.7551/mitpress/9780262033428.003.0007) of [Combinatorial Auctions](https://doi.org/10.7551/mitpress/9780262033428.001.0001).

This work was presented in the following publications: 

**Combinatorial Auction-Based Resource Allocation in the Fog**, *Fawcett, L., Broadbent, M. H. & Race, N. J. P.*, 10/2016 Fifth European Workshop on Software Defined Networks (EWSDN), 2016. IEEE

A more detailed tutorial and description are to follow.

## Install

Airship is built in [`golang`](https://golang.org/).

To fetch Airship, first fetch the respository:

```$ git clone https://github.com/broadbent/airship.git```

Move into the repository:

```$ cd airship```

The run Airship with:

```$ go run airship.go```

### Requirements

Airship has a number of `golang` dependencies. These are handled with [`godep`](https://github.com/tools/godep), so no need for additional installs.

A running [MongoDB](https://www.mongodb.com/) instance is also required. Please see their [documentation](https://www.mongodb.com/download-center?jmp=nav) for instructions on how to install this.

As mentioned previously, Airship works alongside [Siren](https://github.com/lyndon160/Siren-Provisioner) to discover resources and provision successful auctions. A running instance of this is also required. Please see their [documentation](https://github.com/lyndon160/Siren-Provisioner/blob/master/README.md) for more details.

### Configuration

An example default configuration can be found in `airship\config\config.json`, or alternatively [online](https://github.com/broadbent/airship/blob/master/config/config.go).

This contains some sane default configuration. It may be necessary to change this to match your deployment scenario. For example, it may be necessary to change the way Airship connects to the required MongoDB.