import { TypedDocumentNode as DocumentNode } from "@graphql-typed-document-node/core";
export type Maybe<T> = T | null;
export type InputMaybe<T> = Maybe<T>;
export type Exact<T extends { [key: string]: unknown }> = {
  [K in keyof T]: T[K];
};
export type MakeOptional<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]?: Maybe<T[SubKey]>;
};
export type MakeMaybe<T, K extends keyof T> = Omit<T, K> & {
  [SubKey in K]: Maybe<T[SubKey]>;
};
export type MakeEmpty<
  T extends { [key: string]: unknown },
  K extends keyof T,
> = { [_ in K]?: never };
export type Incremental<T> =
  | T
  | {
      [P in keyof T]?: P extends " $fragmentName" | "__typename" ? T[P] : never;
    };
/** All built-in and custom scalars, mapped to their actual values */
export type Scalars = {
  ID: { input: string; output: string };
  String: { input: string; output: string };
  Boolean: { input: boolean; output: boolean };
  Int: { input: number; output: number };
  Float: { input: number; output: number };
  Time: { input: any; output: any };
  join__FieldSet: { input: any; output: any };
  link__Import: { input: any; output: any };
};

export type Aircraft = {
  __typename?: "Aircraft";
  capacity: Scalars["Int"]["output"];
  id: Scalars["ID"]["output"];
  manufacturer: Scalars["String"]["output"];
  model: Scalars["String"]["output"];
  registration: Scalars["String"]["output"];
  status: AircraftStatus;
  yearOfManufacture: Scalars["Int"]["output"];
};

export enum AircraftStatus {
  Available = "AVAILABLE",
  Grounded = "GROUNDED",
  InService = "IN_SERVICE",
  Maintenance = "MAINTENANCE",
}

export type CreateAircraftInput = {
  capacity: Scalars["Int"]["input"];
  manufacturer: Scalars["String"]["input"];
  model: Scalars["String"]["input"];
  registration: Scalars["String"]["input"];
  status: AircraftStatus;
  yearOfManufacture: Scalars["Int"]["input"];
};

export type Flight = {
  __typename?: "Flight";
  aircraft?: Maybe<Aircraft>;
  arrivalTime: Scalars["Time"]["output"];
  departureTime: Scalars["Time"]["output"];
  destination: Scalars["String"]["output"];
  id: Scalars["ID"]["output"];
  number: Scalars["String"]["output"];
  origin: Scalars["String"]["output"];
  status: FlightStatus;
};

export enum FlightStatus {
  Arrived = "ARRIVED",
  Cancelled = "CANCELLED",
  Delayed = "DELAYED",
  Departed = "DEPARTED",
  InProgress = "IN_PROGRESS",
  Scheduled = "SCHEDULED",
  Unspecified = "UNSPECIFIED",
}

export type Mutation = {
  __typename?: "Mutation";
  createAircraft: Aircraft;
  createFlight: Flight;
};

export type MutationCreateAircraftArgs = {
  input?: InputMaybe<CreateAircraftInput>;
};

export type MutationCreateFlightArgs = {
  aircraftId: Scalars["ID"]["input"];
  arrivalTime: Scalars["Time"]["input"];
  departureTime: Scalars["Time"]["input"];
  destination: Scalars["String"]["input"];
  number: Scalars["String"]["input"];
  origin: Scalars["String"]["input"];
};

export type Query = {
  __typename?: "Query";
  getAircraftById?: Maybe<Aircraft>;
  getFlightById?: Maybe<Flight>;
};

export type QueryGetAircraftByIdArgs = {
  input: Scalars["ID"]["input"];
};

export type QueryGetFlightByIdArgs = {
  id: Scalars["ID"]["input"];
};

export enum Join__Graph {
  Aircraft = "AIRCRAFT",
  Flights = "FLIGHTS",
}

export enum Link__Purpose {
  /** `EXECUTION` features provide metadata necessary for operation execution. */
  Execution = "EXECUTION",
  /** `SECURITY` features provide metadata necessary to securely resolve fields. */
  Security = "SECURITY",
}

export type CreateAircraftMutationVariables = Exact<{
  input: CreateAircraftInput;
}>;

export type CreateAircraftMutation = {
  __typename?: "Mutation";
  createAircraft: {
    __typename?: "Aircraft";
    id: string;
    registration: string;
    manufacturer: string;
    model: string;
    status: AircraftStatus;
  };
};

export type CreateFlightMutationVariables = Exact<{
  number: Scalars["String"]["input"];
  origin: Scalars["String"]["input"];
  destination: Scalars["String"]["input"];
  departureTime: Scalars["Time"]["input"];
  arrivalTime: Scalars["Time"]["input"];
  aircraftId: Scalars["ID"]["input"];
}>;

export type CreateFlightMutation = {
  __typename?: "Mutation";
  createFlight: {
    __typename?: "Flight";
    id: string;
    number: string;
    origin: string;
    destination: string;
    departureTime: any;
    arrivalTime: any;
    aircraft?: { __typename?: "Aircraft"; id: string } | null;
  };
};

export type GetAircraftQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

export type GetAircraftQuery = {
  __typename?: "Query";
  getAircraftById?: {
    __typename?: "Aircraft";
    id: string;
    registration: string;
    manufacturer: string;
    model: string;
    status: AircraftStatus;
  } | null;
};

export type GetFlightByIdQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

export type GetFlightByIdQuery = {
  __typename?: "Query";
  getFlightById?: {
    __typename?: "Flight";
    id: string;
    number: string;
    origin: string;
    destination: string;
    status: FlightStatus;
    aircraft?: { __typename?: "Aircraft"; id: string } | null;
  } | null;
};

export type GetFlightByIdWithAircraftQueryVariables = Exact<{
  id: Scalars["ID"]["input"];
}>;

export type GetFlightByIdWithAircraftQuery = {
  __typename?: "Query";
  getFlightById?: {
    __typename?: "Flight";
    id: string;
    number: string;
    origin: string;
    destination: string;
    status: FlightStatus;
    aircraft?: {
      __typename?: "Aircraft";
      id: string;
      registration: string;
      manufacturer: string;
      model: string;
      status: AircraftStatus;
    } | null;
  } | null;
};

export const CreateAircraftDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateAircraft" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "input" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "CreateAircraftInput" },
            },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createAircraft" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "input" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "registration" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "manufacturer" },
                },
                { kind: "Field", name: { kind: "Name", value: "model" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateAircraftMutation,
  CreateAircraftMutationVariables
>;
export const CreateFlightDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "mutation",
      name: { kind: "Name", value: "CreateFlight" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "number" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "origin" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "destination" },
          },
          type: {
            kind: "NonNullType",
            type: {
              kind: "NamedType",
              name: { kind: "Name", value: "String" },
            },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "departureTime" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "Time" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "arrivalTime" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "Time" } },
          },
        },
        {
          kind: "VariableDefinition",
          variable: {
            kind: "Variable",
            name: { kind: "Name", value: "aircraftId" },
          },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "createFlight" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "number" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "number" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "origin" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "origin" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "destination" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "destination" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "departureTime" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "departureTime" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "arrivalTime" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "arrivalTime" },
                },
              },
              {
                kind: "Argument",
                name: { kind: "Name", value: "aircraftId" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "aircraftId" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "number" } },
                { kind: "Field", name: { kind: "Name", value: "origin" } },
                { kind: "Field", name: { kind: "Name", value: "destination" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "departureTime" },
                },
                { kind: "Field", name: { kind: "Name", value: "arrivalTime" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "aircraft" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  CreateFlightMutation,
  CreateFlightMutationVariables
>;
export const GetAircraftDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetAircraft" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getAircraftById" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "input" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "registration" },
                },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "manufacturer" },
                },
                { kind: "Field", name: { kind: "Name", value: "model" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetAircraftQuery, GetAircraftQueryVariables>;
export const GetFlightByIdDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetFlightById" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getFlightById" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "number" } },
                { kind: "Field", name: { kind: "Name", value: "origin" } },
                { kind: "Field", name: { kind: "Name", value: "destination" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "aircraft" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<GetFlightByIdQuery, GetFlightByIdQueryVariables>;
export const GetFlightByIdWithAircraftDocument = {
  kind: "Document",
  definitions: [
    {
      kind: "OperationDefinition",
      operation: "query",
      name: { kind: "Name", value: "GetFlightByIdWithAircraft" },
      variableDefinitions: [
        {
          kind: "VariableDefinition",
          variable: { kind: "Variable", name: { kind: "Name", value: "id" } },
          type: {
            kind: "NonNullType",
            type: { kind: "NamedType", name: { kind: "Name", value: "ID" } },
          },
        },
      ],
      selectionSet: {
        kind: "SelectionSet",
        selections: [
          {
            kind: "Field",
            name: { kind: "Name", value: "getFlightById" },
            arguments: [
              {
                kind: "Argument",
                name: { kind: "Name", value: "id" },
                value: {
                  kind: "Variable",
                  name: { kind: "Name", value: "id" },
                },
              },
            ],
            selectionSet: {
              kind: "SelectionSet",
              selections: [
                { kind: "Field", name: { kind: "Name", value: "id" } },
                { kind: "Field", name: { kind: "Name", value: "number" } },
                { kind: "Field", name: { kind: "Name", value: "origin" } },
                { kind: "Field", name: { kind: "Name", value: "destination" } },
                { kind: "Field", name: { kind: "Name", value: "status" } },
                {
                  kind: "Field",
                  name: { kind: "Name", value: "aircraft" },
                  selectionSet: {
                    kind: "SelectionSet",
                    selections: [
                      { kind: "Field", name: { kind: "Name", value: "id" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "registration" },
                      },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "manufacturer" },
                      },
                      { kind: "Field", name: { kind: "Name", value: "model" } },
                      {
                        kind: "Field",
                        name: { kind: "Name", value: "status" },
                      },
                    ],
                  },
                },
              ],
            },
          },
        ],
      },
    },
  ],
} as unknown as DocumentNode<
  GetFlightByIdWithAircraftQuery,
  GetFlightByIdWithAircraftQueryVariables
>;
