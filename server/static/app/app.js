'use strict';

angular.module('aptWebApp', [
    'ngCookies',
    'ngResource',
    'ngSanitize',
    'ui.router',
    'validation.match',
    'ngFileUpload',
    'ngAnimate',
    'bsTable'
  ])
  .config(function($urlRouterProvider, $locationProvider) {
    $urlRouterProvider
      .otherwise('/table');

    $locationProvider.html5Mode(true);
  });
