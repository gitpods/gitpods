part of swagger.api;



class RepositoriesApi {
  final ApiClient apiClient;

  RepositoriesApi([ApiClient apiClient]) : apiClient = apiClient ?? defaultApiClient;

  /// Create a new repository
  ///
  /// 
  Future<Repository> createRepository(NewRepository newRepository) async {
    Object postBody = newRepository;

    // verify required params are set
    if(newRepository == null) {
     throw new ApiException(400, "Missing required param: newRepository");
    }

    // create path and map variables
    String path = "/repositories".replaceAll("{format}","json");

    // query params
    List<QueryParam> queryParams = [];
    Map<String, String> headerParams = {};
    Map<String, String> formParams = {};
    
    List<String> contentTypes = ["application/json"];

    String contentType = contentTypes.length > 0 ? contentTypes[0] : "application/json";
    List<String> authNames = [];

    if(contentType.startsWith("multipart/form-data")) {
      bool hasFields = false;
      MultipartRequest mp = new MultipartRequest(null, null);
      
      if(hasFields)
        postBody = mp;
    }
    else {
          }

    var response = await apiClient.invokeAPI(path,
                                             'POST',
                                             queryParams,
                                             postBody,
                                             headerParams,
                                             formParams,
                                             contentType,
                                             authNames);

    if(response.statusCode >= 400) {
      throw new ApiException(response.statusCode, response.body);
    } else if(response.body != null) {
      return 
          apiClient.deserialize(response.body, 'Repository') as Repository ;
    } else {
      return null;
    }
  }
  /// Get a owner&#39;s repositories
  ///
  /// 
  Future<List<Repository>> getOwnerRepositories(String owner) async {
    Object postBody = null;

    // verify required params are set
    if(owner == null) {
     throw new ApiException(400, "Missing required param: owner");
    }

    // create path and map variables
    String path = "/repositories/{owner}".replaceAll("{format}","json").replaceAll("{" + "owner" + "}", owner.toString());

    // query params
    List<QueryParam> queryParams = [];
    Map<String, String> headerParams = {};
    Map<String, String> formParams = {};
    
    List<String> contentTypes = ["application/json"];

    String contentType = contentTypes.length > 0 ? contentTypes[0] : "application/json";
    List<String> authNames = [];

    if(contentType.startsWith("multipart/form-data")) {
      bool hasFields = false;
      MultipartRequest mp = new MultipartRequest(null, null);
      
      if(hasFields)
        postBody = mp;
    }
    else {
          }

    var response = await apiClient.invokeAPI(path,
                                             'GET',
                                             queryParams,
                                             postBody,
                                             headerParams,
                                             formParams,
                                             contentType,
                                             authNames);

    if(response.statusCode >= 400) {
      throw new ApiException(response.statusCode, response.body);
    } else if(response.body != null) {
      return 
        (apiClient.deserialize(response.body, 'List<Repository>') as List).map((item) => item as Repository).toList();
    } else {
      return null;
    }
  }
  /// Get a repository by owner name and its name
  ///
  /// 
  Future<Repository> getRepository(String owner, String name) async {
    Object postBody = null;

    // verify required params are set
    if(owner == null) {
     throw new ApiException(400, "Missing required param: owner");
    }
    if(name == null) {
     throw new ApiException(400, "Missing required param: name");
    }

    // create path and map variables
    String path = "/repositories/{owner}/{name}".replaceAll("{format}","json").replaceAll("{" + "owner" + "}", owner.toString()).replaceAll("{" + "name" + "}", name.toString());

    // query params
    List<QueryParam> queryParams = [];
    Map<String, String> headerParams = {};
    Map<String, String> formParams = {};
    
    List<String> contentTypes = ["application/json"];

    String contentType = contentTypes.length > 0 ? contentTypes[0] : "application/json";
    List<String> authNames = [];

    if(contentType.startsWith("multipart/form-data")) {
      bool hasFields = false;
      MultipartRequest mp = new MultipartRequest(null, null);
      
      if(hasFields)
        postBody = mp;
    }
    else {
          }

    var response = await apiClient.invokeAPI(path,
                                             'GET',
                                             queryParams,
                                             postBody,
                                             headerParams,
                                             formParams,
                                             contentType,
                                             authNames);

    if(response.statusCode >= 400) {
      throw new ApiException(response.statusCode, response.body);
    } else if(response.body != null) {
      return 
          apiClient.deserialize(response.body, 'Repository') as Repository ;
    } else {
      return null;
    }
  }
  /// Get all branches of a repository
  ///
  /// 
  Future<List<Branch>> getRepositoryBranches(String owner, String name) async {
    Object postBody = null;

    // verify required params are set
    if(owner == null) {
     throw new ApiException(400, "Missing required param: owner");
    }
    if(name == null) {
     throw new ApiException(400, "Missing required param: name");
    }

    // create path and map variables
    String path = "/repositories/{owner}/{name}/branches".replaceAll("{format}","json").replaceAll("{" + "owner" + "}", owner.toString()).replaceAll("{" + "name" + "}", name.toString());

    // query params
    List<QueryParam> queryParams = [];
    Map<String, String> headerParams = {};
    Map<String, String> formParams = {};
    
    List<String> contentTypes = ["application/json"];

    String contentType = contentTypes.length > 0 ? contentTypes[0] : "application/json";
    List<String> authNames = [];

    if(contentType.startsWith("multipart/form-data")) {
      bool hasFields = false;
      MultipartRequest mp = new MultipartRequest(null, null);
      
      if(hasFields)
        postBody = mp;
    }
    else {
          }

    var response = await apiClient.invokeAPI(path,
                                             'GET',
                                             queryParams,
                                             postBody,
                                             headerParams,
                                             formParams,
                                             contentType,
                                             authNames);

    if(response.statusCode >= 400) {
      throw new ApiException(response.statusCode, response.body);
    } else if(response.body != null) {
      return 
        (apiClient.deserialize(response.body, 'List<Branch>') as List).map((item) => item as Branch).toList();
    } else {
      return null;
    }
  }
  /// Get the commits for a repository
  ///
  /// 
  Future<List<Commit>> getRepositoryCommits(String owner, String name, { String ref }) async {
    Object postBody = null;

    // verify required params are set
    if(owner == null) {
     throw new ApiException(400, "Missing required param: owner");
    }
    if(name == null) {
     throw new ApiException(400, "Missing required param: name");
    }

    // create path and map variables
    String path = "/repositories/{owner}/{name}/commits".replaceAll("{format}","json").replaceAll("{" + "owner" + "}", owner.toString()).replaceAll("{" + "name" + "}", name.toString());

    // query params
    List<QueryParam> queryParams = [];
    Map<String, String> headerParams = {};
    Map<String, String> formParams = {};
    if(ref != null) {
      queryParams.addAll(_convertParametersForCollectionFormat("", "ref", ref));
    }
    
    List<String> contentTypes = ["application/json"];

    String contentType = contentTypes.length > 0 ? contentTypes[0] : "application/json";
    List<String> authNames = [];

    if(contentType.startsWith("multipart/form-data")) {
      bool hasFields = false;
      MultipartRequest mp = new MultipartRequest(null, null);
      
      if(hasFields)
        postBody = mp;
    }
    else {
          }

    var response = await apiClient.invokeAPI(path,
                                             'GET',
                                             queryParams,
                                             postBody,
                                             headerParams,
                                             formParams,
                                             contentType,
                                             authNames);

    if(response.statusCode >= 400) {
      throw new ApiException(response.statusCode, response.body);
    } else if(response.body != null) {
      return 
        (apiClient.deserialize(response.body, 'List<Commit>') as List).map((item) => item as Commit).toList();
    } else {
      return null;
    }
  }
  /// Get the tree including folders (tree) and files (blob) for a repository
  ///
  /// 
  Future<List<TreeEntry>> getRepositoryTree(String owner, String name, { String ref, String path }) async {
    Object postBody = null;

    // verify required params are set
    if(owner == null) {
     throw new ApiException(400, "Missing required param: owner");
    }
    if(name == null) {
     throw new ApiException(400, "Missing required param: name");
    }

    // create path and map variables
    String path = "/repositories/{owner}/{name}/tree".replaceAll("{format}","json").replaceAll("{" + "owner" + "}", owner.toString()).replaceAll("{" + "name" + "}", name.toString());

    // query params
    List<QueryParam> queryParams = [];
    Map<String, String> headerParams = {};
    Map<String, String> formParams = {};
    if(ref != null) {
      queryParams.addAll(_convertParametersForCollectionFormat("", "ref", ref));
    }
    if(path != null) {
      queryParams.addAll(_convertParametersForCollectionFormat("", "path", path));
    }
    
    List<String> contentTypes = ["application/json"];

    String contentType = contentTypes.length > 0 ? contentTypes[0] : "application/json";
    List<String> authNames = [];

    if(contentType.startsWith("multipart/form-data")) {
      bool hasFields = false;
      MultipartRequest mp = new MultipartRequest(null, null);
      
      if(hasFields)
        postBody = mp;
    }
    else {
          }

    var response = await apiClient.invokeAPI(path,
                                             'GET',
                                             queryParams,
                                             postBody,
                                             headerParams,
                                             formParams,
                                             contentType,
                                             authNames);

    if(response.statusCode >= 400) {
      throw new ApiException(response.statusCode, response.body);
    } else if(response.body != null) {
      return 
        (apiClient.deserialize(response.body, 'List<TreeEntry>') as List).map((item) => item as TreeEntry).toList();
    } else {
      return null;
    }
  }
}
