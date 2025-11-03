package search.config;

import com.fasterxml.jackson.databind.ObjectMapper;
import com.fasterxml.jackson.datatype.jsr310.JavaTimeModule;
import javax.net.ssl.SSLContext;
import org.apache.http.HttpHost;
import org.apache.http.auth.AuthScope;
import org.apache.http.auth.UsernamePasswordCredentials;
import org.apache.http.client.CredentialsProvider;
import org.apache.http.impl.client.BasicCredentialsProvider;
import org.apache.http.ssl.SSLContexts;
import org.opensearch.client.RestClient;
import org.opensearch.client.RestClientBuilder;
import org.opensearch.client.RestHighLevelClient;
import org.opensearch.client.json.jackson.JacksonJsonpMapper;
import org.opensearch.client.opensearch.OpenSearchClient;
import org.opensearch.client.transport.OpenSearchTransport;
import org.opensearch.client.transport.rest_client.RestClientTransport;
import org.springframework.beans.factory.annotation.Value;
import org.springframework.context.annotation.Bean;
import org.springframework.context.annotation.Configuration;

/**
 * This class contains the opensearch config.
 */
@Configuration
public class OpenSearchConfig {

  @Value("${opensearch.host}")
  private String host;

  @Value("${opensearch.port}")
  private int port;

  @Value("${opensearch.username}")
  private String username;

  @Value("${opensearch.password}")
  private String password;

  @Value("${opensearch.scheme}")
  private String scheme;

  /**
   * This creates an opensearch rest high-level client.
   *
   * @return the rest client.
   */
  @Bean
  public RestHighLevelClient restHighLevelClient() {
    final CredentialsProvider credentialsProvider = new BasicCredentialsProvider();
    credentialsProvider.setCredentials(AuthScope.ANY,
            new UsernamePasswordCredentials(username, password));

    RestClientBuilder builder = RestClient.builder(new HttpHost(host, port, scheme))
            .setHttpClientConfigCallback(httpClientBuilder -> {
              httpClientBuilder.setDefaultCredentialsProvider(credentialsProvider);

              if ("https".equals(scheme)) {
                try {
                  SSLContext sslContext = SSLContexts.custom()
                          .loadTrustMaterial(null, (x509Certificates, s) -> true)
                          .build();
                  httpClientBuilder.setSSLContext(sslContext);
                  httpClientBuilder.setSSLHostnameVerifier((s, sslSession) -> true);
                } catch (Exception e) {
                  throw new RuntimeException("Failed to create SSL context", e);
                }
              }

              return httpClientBuilder;
            });

    return new RestHighLevelClient(builder);
  }

  /**
   * This creates the opensearch client.
   *
   * @param restHighLevelClient an opensearch rest client to use.
   * @return the opensearch client.
   */
  @Bean
  public OpenSearchClient openSearchClient(RestHighLevelClient restHighLevelClient) {
    ObjectMapper objectMapper = new ObjectMapper();
    objectMapper.registerModule(new JavaTimeModule());

    OpenSearchTransport transport = new RestClientTransport(
            restHighLevelClient.getLowLevelClient(),
            new JacksonJsonpMapper(objectMapper)
    );
    return new OpenSearchClient(transport);
  }
}