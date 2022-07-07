// Instance describes a Dapr sidecar instance information
export interface Instance {
    appID: string;
    httpPort: number;
    grpcPort: number;
    appPort: number;
    command: string;
    age: string;
    created: string;
    pid: number;
    replicas: number;
    address: string;
    supportsDeletion: boolean;
    supportsLogs: boolean;
    manifest: string;
    status: string;
    labels: string;
    selector: string;
    config: string;
}

// Status represents the status of a named Dapr resource
export interface Status {
    service: string;
    name: string;
    namespace: string;
    healthy: string;
    status: string;
    version: string;
    age: string;
    created: string;
}

// Metadata represents metadata from dapr sidecar.
export interface Metadata {
    id: string;
    actors: MetadataActors[];
    extended: {[key: string]: any};
}

// MetadataActors represents actor metadata: type and count
export interface MetadataActors {
    type: string;
    count: number;
}

// Log represents a log object supporting log metadata
export interface Log {
    level: string;
    timestamp: number;
    container: string;
    content: string;
}

// DaprComponent describes an Dapr component type
export interface DaprComponent {
    name: string;
    kind: string;
    type: string;
    created: string;
    age: string;
    scopes: string[];
    manifest: any;
    img: string;
}

// DaprConfiguration represents a Dapr configuration
export interface DaprConfiguration {
    name: string;
    kind: string;
    created: string;
    age: string;
    tracingEnabled: boolean;
    samplingRate: string;
    metricsEnabled: boolean;
    mtlsEnabled: boolean;
    mtlsWorkloadTTL: string;
    mtlsClockSkew: string;
    manifest: any;
}

// Dapr version
export interface DaprVersion {
  version: string;
  runtimeVersion: string;
}

