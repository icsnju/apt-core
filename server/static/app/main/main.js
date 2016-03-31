'use strict';

angular.module('aptWebApp')
  .config(function($stateProvider) {
    $stateProvider
      .state('main', {
        templateUrl: 'static/app/main/main.html',
      });
  });
