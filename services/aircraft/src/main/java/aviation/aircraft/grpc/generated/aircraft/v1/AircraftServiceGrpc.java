package generated.aircraft.v1;

import static io.grpc.MethodDescriptor.generateFullMethodName;

/**
 */
@javax.annotation.Generated(
    value = "by gRPC proto compiler (version 1.70.0)",
    comments = "Source: aircraft/v1/aircraft.proto")
@io.grpc.stub.annotations.GrpcGenerated
public final class AircraftServiceGrpc {

  private AircraftServiceGrpc() {}

  public static final java.lang.String SERVICE_NAME = "aircraft.v1.AircraftService";

  // Static method descriptors that strictly reflect the proto.
  private static volatile io.grpc.MethodDescriptor<generated.aircraft.v1.GetAircraftByIdRequest,
      generated.aircraft.v1.GetAircraftByIdResponse> getGetAircraftByIdMethod;

  @io.grpc.stub.annotations.RpcMethod(
      fullMethodName = SERVICE_NAME + '/' + "GetAircraftById",
      requestType = generated.aircraft.v1.GetAircraftByIdRequest.class,
      responseType = generated.aircraft.v1.GetAircraftByIdResponse.class,
      methodType = io.grpc.MethodDescriptor.MethodType.UNARY)
  public static io.grpc.MethodDescriptor<generated.aircraft.v1.GetAircraftByIdRequest,
      generated.aircraft.v1.GetAircraftByIdResponse> getGetAircraftByIdMethod() {
    io.grpc.MethodDescriptor<generated.aircraft.v1.GetAircraftByIdRequest, generated.aircraft.v1.GetAircraftByIdResponse> getGetAircraftByIdMethod;
    if ((getGetAircraftByIdMethod = AircraftServiceGrpc.getGetAircraftByIdMethod) == null) {
      synchronized (AircraftServiceGrpc.class) {
        if ((getGetAircraftByIdMethod = AircraftServiceGrpc.getGetAircraftByIdMethod) == null) {
          AircraftServiceGrpc.getGetAircraftByIdMethod = getGetAircraftByIdMethod =
              io.grpc.MethodDescriptor.<generated.aircraft.v1.GetAircraftByIdRequest, generated.aircraft.v1.GetAircraftByIdResponse>newBuilder()
              .setType(io.grpc.MethodDescriptor.MethodType.UNARY)
              .setFullMethodName(generateFullMethodName(SERVICE_NAME, "GetAircraftById"))
              .setSampledToLocalTracing(true)
              .setRequestMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  generated.aircraft.v1.GetAircraftByIdRequest.getDefaultInstance()))
              .setResponseMarshaller(io.grpc.protobuf.ProtoUtils.marshaller(
                  generated.aircraft.v1.GetAircraftByIdResponse.getDefaultInstance()))
              .setSchemaDescriptor(new AircraftServiceMethodDescriptorSupplier("GetAircraftById"))
              .build();
        }
      }
    }
    return getGetAircraftByIdMethod;
  }

  /**
   * Creates a new async stub that supports all call types for the service
   */
  public static AircraftServiceStub newStub(io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AircraftServiceStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AircraftServiceStub>() {
        @java.lang.Override
        public AircraftServiceStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AircraftServiceStub(channel, callOptions);
        }
      };
    return AircraftServiceStub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports all types of calls on the service
   */
  public static AircraftServiceBlockingV2Stub newBlockingV2Stub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AircraftServiceBlockingV2Stub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AircraftServiceBlockingV2Stub>() {
        @java.lang.Override
        public AircraftServiceBlockingV2Stub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AircraftServiceBlockingV2Stub(channel, callOptions);
        }
      };
    return AircraftServiceBlockingV2Stub.newStub(factory, channel);
  }

  /**
   * Creates a new blocking-style stub that supports unary and streaming output calls on the service
   */
  public static AircraftServiceBlockingStub newBlockingStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AircraftServiceBlockingStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AircraftServiceBlockingStub>() {
        @java.lang.Override
        public AircraftServiceBlockingStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AircraftServiceBlockingStub(channel, callOptions);
        }
      };
    return AircraftServiceBlockingStub.newStub(factory, channel);
  }

  /**
   * Creates a new ListenableFuture-style stub that supports unary calls on the service
   */
  public static AircraftServiceFutureStub newFutureStub(
      io.grpc.Channel channel) {
    io.grpc.stub.AbstractStub.StubFactory<AircraftServiceFutureStub> factory =
      new io.grpc.stub.AbstractStub.StubFactory<AircraftServiceFutureStub>() {
        @java.lang.Override
        public AircraftServiceFutureStub newStub(io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
          return new AircraftServiceFutureStub(channel, callOptions);
        }
      };
    return AircraftServiceFutureStub.newStub(factory, channel);
  }

  /**
   */
  public interface AsyncService {

    /**
     */
    default void getAircraftById(generated.aircraft.v1.GetAircraftByIdRequest request,
        io.grpc.stub.StreamObserver<generated.aircraft.v1.GetAircraftByIdResponse> responseObserver) {
      io.grpc.stub.ServerCalls.asyncUnimplementedUnaryCall(getGetAircraftByIdMethod(), responseObserver);
    }
  }

  /**
   * Base class for the server implementation of the service AircraftService.
   */
  public static abstract class AircraftServiceImplBase
      implements io.grpc.BindableService, AsyncService {

    @java.lang.Override public final io.grpc.ServerServiceDefinition bindService() {
      return AircraftServiceGrpc.bindService(this);
    }
  }

  /**
   * A stub to allow clients to do asynchronous rpc calls to service AircraftService.
   */
  public static final class AircraftServiceStub
      extends io.grpc.stub.AbstractAsyncStub<AircraftServiceStub> {
    private AircraftServiceStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AircraftServiceStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AircraftServiceStub(channel, callOptions);
    }

    /**
     */
    public void getAircraftById(generated.aircraft.v1.GetAircraftByIdRequest request,
        io.grpc.stub.StreamObserver<generated.aircraft.v1.GetAircraftByIdResponse> responseObserver) {
      io.grpc.stub.ClientCalls.asyncUnaryCall(
          getChannel().newCall(getGetAircraftByIdMethod(), getCallOptions()), request, responseObserver);
    }
  }

  /**
   * A stub to allow clients to do synchronous rpc calls to service AircraftService.
   */
  public static final class AircraftServiceBlockingV2Stub
      extends io.grpc.stub.AbstractBlockingStub<AircraftServiceBlockingV2Stub> {
    private AircraftServiceBlockingV2Stub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AircraftServiceBlockingV2Stub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AircraftServiceBlockingV2Stub(channel, callOptions);
    }

    /**
     */
    public generated.aircraft.v1.GetAircraftByIdResponse getAircraftById(generated.aircraft.v1.GetAircraftByIdRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAircraftByIdMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do limited synchronous rpc calls to service AircraftService.
   */
  public static final class AircraftServiceBlockingStub
      extends io.grpc.stub.AbstractBlockingStub<AircraftServiceBlockingStub> {
    private AircraftServiceBlockingStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AircraftServiceBlockingStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AircraftServiceBlockingStub(channel, callOptions);
    }

    /**
     */
    public generated.aircraft.v1.GetAircraftByIdResponse getAircraftById(generated.aircraft.v1.GetAircraftByIdRequest request) {
      return io.grpc.stub.ClientCalls.blockingUnaryCall(
          getChannel(), getGetAircraftByIdMethod(), getCallOptions(), request);
    }
  }

  /**
   * A stub to allow clients to do ListenableFuture-style rpc calls to service AircraftService.
   */
  public static final class AircraftServiceFutureStub
      extends io.grpc.stub.AbstractFutureStub<AircraftServiceFutureStub> {
    private AircraftServiceFutureStub(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      super(channel, callOptions);
    }

    @java.lang.Override
    protected AircraftServiceFutureStub build(
        io.grpc.Channel channel, io.grpc.CallOptions callOptions) {
      return new AircraftServiceFutureStub(channel, callOptions);
    }

    /**
     */
    public com.google.common.util.concurrent.ListenableFuture<generated.aircraft.v1.GetAircraftByIdResponse> getAircraftById(
        generated.aircraft.v1.GetAircraftByIdRequest request) {
      return io.grpc.stub.ClientCalls.futureUnaryCall(
          getChannel().newCall(getGetAircraftByIdMethod(), getCallOptions()), request);
    }
  }

  private static final int METHODID_GET_AIRCRAFT_BY_ID = 0;

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
        case METHODID_GET_AIRCRAFT_BY_ID:
          serviceImpl.getAircraftById((generated.aircraft.v1.GetAircraftByIdRequest) request,
              (io.grpc.stub.StreamObserver<generated.aircraft.v1.GetAircraftByIdResponse>) responseObserver);
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
          getGetAircraftByIdMethod(),
          io.grpc.stub.ServerCalls.asyncUnaryCall(
            new MethodHandlers<
              generated.aircraft.v1.GetAircraftByIdRequest,
              generated.aircraft.v1.GetAircraftByIdResponse>(
                service, METHODID_GET_AIRCRAFT_BY_ID)))
        .build();
  }

  private static abstract class AircraftServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoFileDescriptorSupplier, io.grpc.protobuf.ProtoServiceDescriptorSupplier {
    AircraftServiceBaseDescriptorSupplier() {}

    @java.lang.Override
    public com.google.protobuf.Descriptors.FileDescriptor getFileDescriptor() {
      return generated.aircraft.v1.AircraftProto.getDescriptor();
    }

    @java.lang.Override
    public com.google.protobuf.Descriptors.ServiceDescriptor getServiceDescriptor() {
      return getFileDescriptor().findServiceByName("AircraftService");
    }
  }

  private static final class AircraftServiceFileDescriptorSupplier
      extends AircraftServiceBaseDescriptorSupplier {
    AircraftServiceFileDescriptorSupplier() {}
  }

  private static final class AircraftServiceMethodDescriptorSupplier
      extends AircraftServiceBaseDescriptorSupplier
      implements io.grpc.protobuf.ProtoMethodDescriptorSupplier {
    private final java.lang.String methodName;

    AircraftServiceMethodDescriptorSupplier(java.lang.String methodName) {
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
      synchronized (AircraftServiceGrpc.class) {
        result = serviceDescriptor;
        if (result == null) {
          serviceDescriptor = result = io.grpc.ServiceDescriptor.newBuilder(SERVICE_NAME)
              .setSchemaDescriptor(new AircraftServiceFileDescriptorSupplier())
              .addMethod(getGetAircraftByIdMethod())
              .build();
        }
      }
    }
    return result;
  }
}
