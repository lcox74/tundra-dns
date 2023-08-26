
export const RecordType = {
    A: 1,
    NS: 2,
    CNAME: 5,
    MX: 15,
    TXT: 16,
};

export abstract class Record {
    id: number;
    domain: string;
    subdomain: string;
    
    type: number;
    routeType: number;

    createdAt: Date;
    updatedAt: Date;
    deactivatedAt: Date;
    expiredAt: Date;

    ttl: number;

    allow: string[];
    deny: string[];

    constructor(data: any) {
        this.id = data.id;
        this.domain = data.domain;
        this.subdomain = data.subdomain;
        this.type = data.type;
        this.routeType = data.route_type;
        this.createdAt = new Date(data.created_at);
        this.updatedAt = new Date(data.updated_at);
        this.deactivatedAt = new Date(data.deactivated_at);
        this.expiredAt = new Date(data.expired_at);
        this.ttl = data.ttl;
        this.allow = data.allow;
        this.deny = data.deny;
    }
}

export class ARecord extends Record {
    address: string;
    lastSeen: Date;

    constructor(data: any) {
        super(data);
        this.address = data.data.address;
        this.lastSeen = new Date(data.data.last_seen);
    }
}

export class CNAMERecord extends Record {
    alias: string;

    constructor(data: any) {
        super(data);
        this.alias = data.data.alias;
    }
}

