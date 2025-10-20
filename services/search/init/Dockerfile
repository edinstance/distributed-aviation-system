FROM alpine:3.18

# Install curl and bash for API calls
RUN apk add --no-cache curl bash jq

# Copy the setup script
COPY setup-dashboard.sh /usr/local/bin/setup-dashboard.sh
RUN chmod +x /usr/local/bin/setup-dashboard.sh

# Copy all the init files
COPY indexes/ /init/indexes/
COPY visualizations/ /init/visualizations/
COPY searches/ /init/searches/
COPY dashboards/ /init/dashboards/

# Set the default command
CMD ["/usr/local/bin/setup-dashboard.sh"]