# Plan
> 25th Aug 2023

TundraDNS is a comprehensive and user-friendly platform designed to simplify 
the management of domain names and their associated DNS records. The following
is a spec, outlining the architecture, features and technical details of 
TundraDNS. 

## Scope

TundraDNS aims to streamline the process of managing domain names and DNS 
record for individual devs who want to experiment with IoT devices, to small 
businesses wanting to automatically manage ephemeral services. The service will 
incorporate a range of routing policies, health checks, and failover mechanisms 
to ensure high availability and performance of resources.

## High-Level Architecture

The TundraDNS Service is built on a modular and scalable architecture to ensure 
high availability, reliability, and ease of maintenance. The architecture 
consists of the following key components:

- **User Interface**: The web or cli based user interface which provides the 
  user with a seamless experience for managing domains, DNS records and routing 
  configurations.
- **API Gateway**: The API gateway acts as a bridge between the user interface 
  and the backend services, facilitating secure communication and data exchange.
- **Routing Engine**: The routing engine is responsible for processing routing 
  policies, health checks, and failover mechanisms to optimise traffic 
  distribution.
- **Cache**: The cache stores frequently accessed DNS records, reducing the 
  need to repeatedly query other authoritative DNS servers or the Database and 
  improving response times.
- **Database**: The database stores user account information, domain details, 
  DNS records, and configuration settings.
- **Authentication and Security Layer**: This component ensures a secure user 
  or machine authentication, authorisation and data encryption throughout the 
  platform.

## Architecture

Though the high-level architecture has been described, time to get into the 
details on the system. This will include how TundraDNS will handle each of the
elements described above as well as how the elements communicate with each 
other.

<p align="center">
  <img src="./res/architecture-darkmode.svg#gh-dark-mode-only" alt="Light Mode" width="400"/>
  <img src="./res/architecture-lightmode.svg#gh-light-mode-only" alt="Dark Mode" width="400"/>
</p>

### Authentication

In terms of the actual architecture, the biggest pain will be Authentication. 
I do have a plan for this though. Dont do authentication. Instead, because this
is designed to be a self hosted application the management web interface can be
"public" on the cavet that its only exposed on an internal network. For this
we can setup a Tailscale virtual network which the webapp and API will be 
exposed on.

The only traffic that wont be bound to the Tailscale network will be all DNS
queries which will come in on port `:53/udp`. I don't have plans on dealing with
dns over TLS/TCP (though the library I'm using says it supports it).

### User Interface

The user interface will be a webapp which will be exposed on the Tailscale 
network as described in the Authentication section. This will require displaying
the records and managing them accordinly. For the MVP there will only be the 
option of a single Domain to manage, though on the backend it should technically
still be able to handle more. The following will be used to build the webapp:

- Vue (The Framework)
- Vite (The builder/packager)
- Pinia (The state management)
- Tailwindcss (The I'm too lazy to do my own CSS)

### Backend

The backend for this project is designed to be a single binary. It's built up of
3 distinct modules: Routing Engine, Record Management API and DNS Query Handler.
These will function seperately but do share some form of communication. This 
will also be written in [Go](https://go.dev/) for its simplicity, speed and 
management of threading.

**Routing Engine**

The Routing Engine is the core system that defines and manages everything. It 
receives messages from the Record Management API and performs the tasks to make
it happen. Persistent storage of records will be put in the main database which
the API can read.

The Routing Engine will periodically build/sync/upate a full routing table in 
the Redis Cache. This is so the DNS Resolver has quick access to the active 
records without needing to do any additional processing.

**API Gateway**

The API Gateway has Read access to the main Database, I don't want the external
facing API to directly affect the DB if I can handle it. So when a request to 
add, update or delete gets made, then it must notify the Routing Engine with the
request data so it can make the enformed decision on what to do. This will not
happen cocurrently as the appropriate response needs to be sent back to the 
client.

**DNS Query Handler**

As mentioned previously, this is designed to be a simple module. It recieves 
DNS requests, looks through the request's `DNS Questions` and then polls the 
Redis Cache for the appropriate answer. There should be little to no additional
processing done (except for weighted load balancers).

### Database and Cache

The main Database will be persistent storage of the full record data. This will
include additional meta data such as loadbalancing information, what nodes are 
active, is the record still valid etc. The Routing Engine processes all this 
data and builds a stripped down and pre-processed record data so the DNS Query
Handler can quickly respond to requests. 

| Column Name      | Data Type | Description                                   |
| ---------------- | --------- | ----------------------------------------------|
| `id`             | INTEGER   | Primary Key                                   |
| `domain`         | TEXT      | The base domain of the record                 |
| `subdomain`      | TEXT      | The subdomain, added to the domain = FQDN     |
| `ttl`            | INTEGER   | Cache TTL for the record in seconds           |
| `type`           | INTEGER   | Record Type (`A`, `CNAME`, `NS` etc.)         |
| `route_type`     | INTEGER   | Custom routing type (`single`, `loadbalance`) |
| `created_at`     | DATETIME  | When the record was added                     |
| `updated_at`     | DATETIME  | Last Updated Time                             |
| `deactivated_at` | DATETIME  | When record was deactivated NULLABLE          |
| `expired_at`     | DATETIME  | If Ephemeral, when to expire NULLABLE         |
| `allow`          | BLOB      | Allow list, comma seperated IP/Subnet NULLABLE|
| `deny`           | BLOB      | Block list, comma seperated IP/Subnet NULLABLE|
| `data`           | BLOB      | JSON data, specific to record type            |

**A RECORD**

The data structure of the `A Record` is as follows:

```json
{
    "address": "10.10.10.10",   // Target Address
    "last_seen": "",            // Timestamp (0001-01-01T00:00:00Z)
    "fallback": "10.10.10.11"   // Fallback IP (Optional)
}
```

### Logging
> **Note:** This will not be developed in time, probably.

The Logging aspect of TundraDNS is simple, but runs as a seperate entity. The 
backend service logs everything in a structurely manner into `stdout`, this can
be captured using [vector](https://vector.dev) and automatically be directed to
a Database (Like InfluxDB). From this data a Grafana instance can pull the logs
from the DB and generate analyitics for us.

## Frontend Design

By default, going to the webapp should display the Home screen. This will 
display all records for the domain as a minimal row based table. Though simple
it should be interactive and intuitive.

![Home Screen](./res/frontend/Home.png)

The border of rows must highlight on hover to indicate possibility of selection.
Then when selected, it should expand the row downwards showing a more detailed
view of the record.

![Selected Record](./res/frontend/Select%20Single%20Record.png)

When creating a new record, a popup form should be displayed which changes 
dynamically based on the selected type. This popup form will be used for 
updating records as well with the button text changing accordingly.

![Create/Update Record](res/frontend/Create%20Record.png)

Though the design shows different routing types for the records, the priority is
to handle single A Records before handling others. Deleting a record is also not
currently the priority, though the button is there it wont do anything yet.