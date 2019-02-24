import 'dart:async';
import 'dart:convert';

import 'package:angular/angular.dart';
import 'package:sourcepods/api.dart';
import 'package:sourcepods/src/api/api.dart' as api;
import 'package:sourcepods/src/repository/repository.dart';
import 'package:sourcepods/src/repository/repository_create_component.dart';
import 'package:sourcepods/src/repository/tree.dart';
import 'package:sourcepods/src/repository/commit.dart';

@Injectable()
class RepositoryService {
  RepositoryService(this._api);

  final API _api;

  Future<Repository> get(String owner, String name) async {
    api.Repository r = await _api.repositories.getRepository(owner, name);
    return Repository.fromAPI(r);
  }

  Future<List<String>> getBranches(String owner, String name) async {
    List<api.Branch> branches =
        await _api.repositories.getRepositoryBranches(owner, name);
    return branches.map((b) => b.name).toList();
  }

  Future<List<TreeEntry>> getTree(
      String owner, String name, String ref, String gitPath) async {
    List<api.TreeEntry> tree = await _api.repositories
        .getRepositoryTree(owner, name, ref: ref, path: gitPath);
    return tree.map((te) => TreeEntry.fromAPI(te)).toList();
  }
  
  Future<List<Commit>> getCommits(String owner, String name, String ref) async{
    List<api.Commit> commits = await _api.repositories
        .getRepositoryCommits(owner, name, ref: ref);
    return commits.map((commit) => Commit.fromAPI(commit)).toList();
  }

  Future<Repository> create(
    String name,
    String description,
    String website,
  ) async {
    api.NewRepository newRepository = api.NewRepository();
    newRepository.name = name;
    newRepository.description = description;
    newRepository.website = website;

    try {
      api.Repository apiRepository =
          await _api.repositories.createRepository(newRepository);
      return Repository.fromAPI(apiRepository);
    } on api.ApiException catch (e) {
      if (e.code == 422) {
        throw new ValidationException.fromJSON(json.decode(e.message));
      } else {
        throw e;
      }
    }
  }
}
