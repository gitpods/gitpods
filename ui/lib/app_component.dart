import 'package:angular/angular.dart';
import 'package:angular_router/angular_router.dart';
import 'package:gitpods/gravatar_component.dart';
import 'package:gitpods/login_component.dart';
import 'package:gitpods/login_service.dart';
import 'package:gitpods/repository_component.dart';
import 'package:gitpods/repository_create_component.dart';
import 'package:gitpods/user.dart';
import 'package:gitpods/user_list_component.dart';
import 'package:gitpods/user_profile_component.dart';
import 'package:gitpods/user_service.dart';
import 'package:gitpods/settings/settings.dart';

@Component(
  selector: 'gitpods-app',
  templateUrl: 'app_component.html',
//  template: '<router-outlet></router-outlet>',
  directives: const [COMMON_DIRECTIVES, ROUTER_DIRECTIVES, Gravatar],
  providers: const [ROUTER_PROVIDERS, UserService, LoginService],
)
@RouteConfig(const [
  const Redirect(
    path: '/',
    redirectTo: const ['UserList'],
  ),
  const Route(
    path: '/users',
    name: 'UserList',
    component: UserListComponent,
  ),
  const Route(
    path: '/login',
    name: 'Login',
    component: LoginComponent,
  ),
  const Route(
    path: '/:username',
    name: 'UserProfile',
    component: UserProfileComponent,
  ),
  const Route(
    path: '/settings/...',
    name: 'Settings',
    component: SettingsComponent,
  ),
  const Route(
    path: '/:owner/:name',
    name: 'Repository',
    component: RepositoryComponent,
  ),
  const Route(
    path: '/new',
    name: 'RepositoryCreate',
    component: RepositoryCreateComponent,
  ),
//  const Route(
//      path: '/**',
//      name: 'NotFound',
//      component: NotFoundComponent,
//  )
])
class AppComponent implements OnInit {
  final Router _router;
  final UserService _userService;
  final LoginService _loginService;

  AppComponent(this._router, this._userService, this._loginService);

  User user;
  bool loading = false;

  @override
  void ngOnInit() {
    this._userService.me()
        .then((User user) => this.user = user)
        .catchError((e) => this._router.navigate(['Login']));
  }

  void logout() {
    this._loginService.logout();
  }
}
