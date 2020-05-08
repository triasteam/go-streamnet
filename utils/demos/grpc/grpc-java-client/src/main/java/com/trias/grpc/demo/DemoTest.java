package com.trias.grpc.demo;

import java.io.UnsupportedEncodingException;
import java.nio.charset.Charset;

import com.google.protobuf.ByteString;
import com.test.rpc.SayHelloServiceGrpc;
import com.test.rpc.Service;

import io.grpc.netty.shaded.io.grpc.netty.NegotiationType;
import io.grpc.netty.shaded.io.grpc.netty.NettyChannelBuilder;

public class DemoTest {

	private static final String host = "127.0.0.1";
	private static final int port = 41005;

	public static void main(String[] args) {
		io.grpc.Channel channel = NettyChannelBuilder.forAddress(host, port).negotiationType(NegotiationType.PLAINTEXT)
				.build();

		Service.SayHelloRequest req = Service.SayHelloRequest.newBuilder()
				.setName(ByteString.copyFrom("倪明", Charset.forName("utf-8"))).build();

		Service.SayHelloResponse result = SayHelloServiceGrpc.newBlockingStub(channel).sayHello(req);
		try {
			System.out.println(result.getResult().toString("utf-8"));
		} catch (UnsupportedEncodingException e) {
			e.printStackTrace();
		}
	}

}
