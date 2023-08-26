
export const RecordType = {
    A: 1,
    NS: 2,
    CNAME: 5,
    SOA: 6,
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

    getFQDN(): string {
        if (this.subdomain == "@") {
            return this.domain.slice(0,-1);
        } else {
            return this.subdomain + "." + this.domain.slice(0,-1);
        }
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

    getAlias(): string {
        return this.alias.slice(0,-1);
    }
}

export class MXRecord extends Record {
    priority: number;
    target: string;

    constructor(data: any) {
        super(data);
        this.priority = data.data.preference;
        this.target = data.data.mail_server;
    }

    getTarget(): string {
        return this.target.slice(0,-1);
    }
}

export class TXTRecord extends Record {
    content: string[];

    constructor(data: any) {
        super(data);
        this.content = data.data.content;
    }

    getContent(): string {
        return this.content.join(" ");
    }

}

export class NSRecord extends Record {
    nameserver: string;

    constructor(data: any) {
        super(data);
        this.nameserver = data.data.server;
    }

    getNameserver(): string {
        return this.nameserver.slice(0,-1);
    }
}

export class SOARecord extends Record {
    primaryNameServer: string;
    responsiblePerson: string;
    serial: number;
    refresh: number;
    retry: number;
    expire: number;
    minimum: number;

    constructor(data: any) {
        super(data);
        this.primaryNameServer = data.data.ns_server;
        this.responsiblePerson = data.data.ns_mailbox;
        this.serial = data.data.serial;
        this.refresh = data.data.refresh;
        this.retry = data.data.retry;
        this.expire = data.data.expire;
        this.minimum = data.data.minttl;
    }

    getNameserver(): string {
        return this.primaryNameServer.slice(0,-1);
    }

    getMailbox(): string {
        let parts = this.responsiblePerson.split(".");

        if (parts.length == 1) {
            return parts[0];
        }

        return parts[0] + '@' + parts.slice(1).join(".");
    }

    getSerial(): string {
        // Convert to Hex

        let hexString = this.serial.toString(16);

        // Pad with leading zeros to ensure 32-bit length
        while (hexString.length < 8) {
            hexString = '0' + hexString;
        }

        return '0x' + hexString;
    }
}