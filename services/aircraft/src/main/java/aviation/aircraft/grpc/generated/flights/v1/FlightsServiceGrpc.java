package generated.flights.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.70.0)",
    comments = "Source: flights/v1/flights.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class FlightsServiceGrpc {

  private FlightsServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "flights.v1.FlightsService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<generated.flights.v1.CreateFlightRequest,
      generated.flights.v1.CreateFlightResponse> getCreateFlightMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "CreateFlight",
      requestType = generated.flights.v1.CreateFlightRequest.class,
      responseType = generated.flights.v1.CreateFlightResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<generated.flights.v1.CreateFlightRequest,
      generated.flights.v1.CreateFlightResponse> getCreateFlightMethod() {
    io.grpc.MethodDescriptor<generated.flights.v1.CreateFlightRequest, generated.flights.v1.CreateFlightResponse> getCreateFlightMethod;
    if ((getCreateFlightMethod = FlightsServiceGrpc.getCreateFlightMethod) == null) {
      synchronized (FlightsServiceGrpc.class) {
        if ((getCreateFlightMethod = FlightsServiceGrpc.getCreateFlightMethod) == null) {
          FlightsServiceGrpc.getCreateFlightMethod = getCreateFlightMethod =
              io.grpc.MethodDescriptor.<generated.flights.v1.CreateFlightRequest, generated.flights.v1.CreateFlightResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "CreateFlight"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  generated.flights.v1.CreateFlightRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  generated.flights.v1.CreateFlightResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlightsServiceMethodDescriptorSupplier("CreateFlight"))
              .build();
        }
      }
    }
    return getCreateFlightMethod;
  }

  private static volatile io.grpc.MethodDescriptor<generated.flights.v1.GetFlightByIdRequest,
      generated.flights.v1.GetFlightByIdResponse> getGetFlightByIdMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetFlightById",
      requestType = generated.flights.v1.GetFlightByIdRequest.class,
      responseType = generated.flights.v1.GetFlightByIdResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<generated.flights.v1.GetFlightByIdRequest,
      generated.flights.v1.GetFlightByIdResponse> getGetFlightByIdMethod() {
    io.grpc.MethodDescriptor<generated.flights.v1.GetFlightByIdRequest, generated.flights.v1.GetFlightByIdResponse> getGetFlightByIdMethod;
    if ((getGetFlightByIdMethod = FlightsServiceGrpc.getGetFlightByIdMethod) == null) {
      synchronized (FlightsServiceGrpc.class) {
        if ((getGetFlightByIdMethod = FlightsServiceGrpc.getGetFlightByIdMethod) == null) {
          FlightsServiceGrpc.getGetFlightByIdMethod = getGetFlightByIdMethod =
              io.grpc.MethodDescriptor.<generated.flights.v1.GetFlightByIdRequest, generated.flights.v1.GetFlightByIdResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetFlightById"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  generated.flights.v1.GetFlightByIdRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  generated.flights.v1.GetFlightByIdResponse.getDefaultInstance()))
              .setSchemaDescriptor(new FlightsServiceMethodDescriptorSupplier("GetFlightById"))
              .build();
        }
      }
    }
    return getGetFlightByIdMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static FlightsServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlightsServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlightsServiceStub>() {
        @java.lang.Override
        public FlightsServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlightsServiceStub(channel, callOptions);
        }
      };
    return FlightsServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static FlightsServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlightsServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlightsServiceBlockingV2Stub>() {
        @java.lang.Override
        public FlightsServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlightsServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return FlightsServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static FlightsServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlightsServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlightsServiceBlockingStub>() {
        @java.lang.Override
        public FlightsServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlightsServiceBlockingStub(channel, callOptions);
        }
      };
    return FlightsServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static FlightsServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<FlightsServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<FlightsServiceFutureStub>() {
        @java.lang.Override
        public FlightsServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new FlightsServiceFutureStub(channel, callOptions);
        }
      };
    return FlightsServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public interface AsyncService {

    /**
     * <pre>
     * CreateFlight requires the following gRPC metadata headers:
     *   - x-user-sub: UUID of the authenticated user
     *   - x-org-id: UUID of the user's organization
     *   - x-org-name: String the user's organization name
     * </pre>
     */
    default void createFlight(generated.flights.v1.CreateFlightRequest request,
        io.grpc.stub.StreamObserver<generated.flights.v1.CreateFlightResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getCreateFlightMethod(), responseObserver);
    }

    /**
     */
    default void getFlightById(generated.flights.v1.GetFlightByIdRequest request,
        io.grpc.stub.StreamObserver<generated.flights.v1.GetFlightByIdResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetFlightByIdMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service FlightsService.
   */
  public static abstract class FlightsServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return FlightsServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service FlightsService.
   */
  public static final class FlightsServiceStub
      extends io.grpc.stub.AbstractAsyncStub<FlightsServiceStub> {
    private FlightsServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlightsServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlightsServiceStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateFlight requires the following gRPC metadata headers:
     *   - x-user-sub: UUID of the authenticated user
     *   - x-org-id: UUID of the user's organization
     *   - x-org-name: String the user's organization name
     * </pre>
     */
    public void createFlight(generated.flights.v1.CreateFlightRequest request,
        io.grpc.stub.StreamObserver<generated.flights.v1.CreateFlightResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getCreateFlightMethod(), getCallOptions()), request, responseObserver);
    }

    /**
     */
    public void getFlightById(generated.flights.v1.GetFlightByIdRequest request,
        io.grpc.stub.StreamObserver<generated.flights.v1.GetFlightByIdResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetFlightByIdMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service FlightsService.
   */
  public static final class FlightsServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<FlightsServiceBlockingV2Stub> {
    private FlightsServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlightsServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlightsServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateFlight requires the following gRPC metadata headers:
     *   - x-user-sub: UUID of the authenticated user
     *   - x-org-id: UUID of the user's organization
     *   - x-org-name: String the user's organization name
     * </pre>
     */
    public generated.flights.v1.CreateFlightResponse createFlight(generated.flights.v1.CreateFlightRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateFlightMethod(), getCallOptions(), request);
    }

    /**
     */
    public generated.flights.v1.GetFlightByIdResponse getFlightById(generated.flights.v1.GetFlightByIdRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetFlightByIdMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service FlightsService.
   */
  public static final class FlightsServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<FlightsServiceBlockingStub> {
    private FlightsServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlightsServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlightsServiceBlockingStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateFlight requires the following gRPC metadata headers:
     *   - x-user-sub: UUID of the authenticated user
     *   - x-org-id: UUID of the user's organization
     *   - x-org-name: String the user's organization name
     * </pre>
     */
    public generated.flights.v1.CreateFlightResponse createFlight(generated.flights.v1.CreateFlightRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getCreateFlightMethod(), getCallOptions(), request);
    }

    /**
     */
    public generated.flights.v1.GetFlightByIdResponse getFlightById(generated.flights.v1.GetFlightByIdRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetFlightByIdMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service FlightsService.
   */
  public static final class FlightsServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<FlightsServiceFutureStub> {
    private FlightsServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected FlightsServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new FlightsServiceFutureStub(channel, callOptions);
    }

    /**
     * <pre>
     * CreateFlight requires the following gRPC metadata headers:
     *   - x-user-sub: UUID of the authenticated user
     *   - x-org-id: UUID of the user's organization
     *   - x-org-name: String the user's organization name
     * </pre>
     */
    public com.google.common.util.concurrent.ListenableFuture<generated.flights.v1.CreateFlightResponse> createFlight(
        generated.flights.v1.CreateFlightRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getCreateFlightMethod(), getCallOptions()), request);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<generated.flights.v1.GetFlightByIdResponse> getFlightById(
        generated.flights.v1.GetFlightByIdRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetFlightByIdMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_CREATE_FLIGHT = 0;
  private static final int METHODID_GET_FLIGHT_BY_ID = 1;

  private static final class MethodHandlers<Req, Resp> implements
      io.grpc.stub.ServerCalls.UnaryMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ServerStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.ClientStreamingMethod<Req, Resp>,
      io.grpc.stub.ServerCalls.BidiStreamingMethod<Req, Resp> {
    private final AsyncService serviceImpl;
    private final int methodId;

    MethodHandlers(AsyncService serviceImpl, int methodId) {
      this.serviceImpl = serviceImpl;
      this.methodId = methodId;
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public void invoke(Req request, io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        case METHODID_CREATE_FLIGHT:
          serviceImpl.createFlight((generated.flights.v1.CreateFlightRequest) request,
              (io.grpc.stub.StreamObserver<generated.flights.v1.CreateFlightResponse>) responseObserver);
          break;
        case METHODID_GET_FLIGHT_BY_ID:
          serviceImpl.getFlightById((generated.flights.v1.GetFlightByIdRequest) request,
              (io.grpc.stub.StreamObserver<generated.flights.v1.GetFlightByIdResponse>) responseObserver);
          break;
        default:
          throw new AssertionError();
      }
    }

    @java.lang.Override
    @java.lang.SuppressWarnings("unchecked")
    public io.grpc.stub.StreamObserver<Req> invoke(
        io.grpc.stub.StreamObserver<Resp> responseObserver) {
      switch (methodId) {
        default:
          throw new AssertionError();
      }
    }
  }

  public static final io.grpc.ServerServiceDefinition bindService(AsyncService service) {
    return io.grpc.ServerServiceDefinition.builder(getServiceDescriptor())
        .addMethod(
          getCreateFlightMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              generated.flights.v1.CreateFlightRequest,
              generated.flights.v1.CreateFlightResponse>(
                service, METHODID_CREATE_FLIGHT)))
        .addMethod(
          getGetFlightByIdMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              generated.flights.v1.GetFlightByIdRequest,
              generated.flights.v1.GetFlightByIdResponse>(
                service, METHODID_GET_FLIGHT_BY_ID)))
        .build();
  }

  private static abstract class FlightsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    FlightsServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return generated.flights.v1.FlightsProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("FlightsService");
    }
  }

  private static final class FlightsServiceFileDescriptorSupplier
      extends FlightsServiceBaseDescriptorSupplier {
    FlightsServiceFileDescriptorSupplier() {}
  }

  private static final class FlightsServiceMethodDescriptorSupplier
      extends FlightsServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    FlightsServiceMethodDescriptorSupplier(java.lang.String methodName) {
      this.methodName = methodName;
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.MethodDescriptor getMethodDescriptor() {
      return getServiceDescriptor().findMethodByName(methodName);
    }
  }

  private static volatile io.grpc.ServiceDescriptor serviceDescriptor;

  public static io.grpc.ServiceDescriptor getServiceDescriptor() {
    io.grpc.ServiceDescriptor result = serviceDescriptor;
    if (result == null) {
      synchronized (FlightsServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new FlightsServiceFileDescriptorSupplier())
              .addMethod(getCreateFlightMethod())
              .addMethod(getGetFlightByIdMethod())
              .build();
        }
      }
    }
    return result;
  }
}
