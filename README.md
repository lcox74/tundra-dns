<p align="center">
  <img src="./res/logo-darkmode.svg#gh-dark-mode-only" alt="Light Mode" width="400"/>
  <img src="./res/logo-lightmode.svg#gh-light-mode-only" alt="Dark Mode" width="400"/>
</p>

> 25th Aug 2023

Welcome to TundraDNS, your DNS management solution built as a project for the 
UQCS 2023 Hackathon. Effortlessly manage records, optimize traffic, and enhance 
your online presence. Join us on this hackathon journey with TundraDNS!

> **Note:** You can view an overview of the project plan [here](./PLAN.md).

**Key Features**

TundraDNS offers a range of powerful features tailored to simplify and enhance the management of DNS records. Though these are designed to be accessed via a WebApp, all the features will be API compatible.

- **Flexible DNS Record Management**: Users can create, update and delete 
  various DNS record types, including `A`, `CNAME`, `MX`, `TXT` and `NS`.
- **Routing Policies and Traffic Optimisation**: TundraDNS supports weighted 
  routing, latency-based routing and IP group based routing. This enables users
  to optimise traffic distribution of their connected services.
- **Health Checks and Failover**: The service includes automated health checks 
  to monitor the availability of resources. In the event of a resource failure, 
  TundraDNS seamlessly redirects traffic to healthy resources, minimising 
  downtime.
- **IP Blocking/Allowlisting**: Enable users to define boundaries and restrict 
  or allow access to resources based on the users' IP addresses.
- **DNS Record Templates**: Enables users to create templates for commonly used 
  DNS records, simplifying the process of adding consistent records.
- **Ephemeral Records**: TundraDNS supports records that have a lease, these 
  records can be periodically pinged as `alive` to extend its lease.
- **Automated Workflows and Webhooks**: Enables users to configure webhooks or 
  event triggers that are fired based on changes to DNS records, allowing for 
  automated workflows and integrations.

> **Note:** This is actually just a wish list. No idea how much I will get done
> I will update later depending on what I'm able to get done.

## Getting Started

TundraDNS is designed to be easily selfhosted. For this reason, the best way to
deploy this is via Docker Compose.

> **TODO:** Actually write the Docker Compose file. Can't currently as I have no
> idea what I'm able to produce and how I'm going to do the Docker Image.