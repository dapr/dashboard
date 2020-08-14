# Dashboard Features

## v0.2.0
* Full integration with Dapr CLI: `dapr dashboard` and `dapr dashboard -k`
* Actors tab added to application detail view
    * View actor count for each type
* Dapr scope filter for applications, components, and configurations
* Components and configurations features added for self-hosted
* Components detail view added
    * Summary: Show summary of component, such as name and type
    * Configuration: Show manifest .yaml file for this component
* Configurations detail view added
    * Summary: Show summary of configuration, such as name and mTLS info
    * (Kubernetes) Applications: Shows and has links to apps currently using this configuration
    * Configuration: Show manifest .yaml file for this configuration
* Light theme/ Dark theme selector
* (Kubernetes) Application logs sorting and filtering
    * User can filter by keyword, container, time since, and date/time range
    * Added functionality to sort logs by ascending and descending date
    * Keyword based highlighting for log messages
    * Added checkboxes to allow user to specify 'only filtered items shown' and 'display timestamps'
* User access to dashboard version in CLI and in the dashboard

## v0.1.0
* Expandable Sidebar with enabled features
* Application table on main page
* (Kubernetes) Control plane 'at-a-glance' on main page
* (Kubernetes) Components view with component name, type, age, and creation timestamp
* (Kubernetes) Configurations view with configuration name, mTLS info, age, and creation timestamp
* (Kubernetes) Control plane tab, mirroring output from `dapr status -k`
* (Standalone) Action to delete running application
* Detail view for applications
    * Summary: instance information (appID, HTTP Port, age, etc.)
    * (Kubernetes) Metadata: Dapr annotations metadata
    * (Kubernetes) Configuration: Deployment manifest from cluster
    * (Kubernetes) Logs: List of daprd logs
        * Filterable by log level: info, debug, warning, error, fatal