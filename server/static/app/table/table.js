'use strict';

angular.module('aptWebApp')
  .config(function($stateProvider) {
    $stateProvider
      .state('main.jobs', {
        url: '/jobs',
        templateUrl: 'static/app/table/jobs.html',
        controller: 'JobCtrl'
      })
      .state('main.jobDetail', {
        url: '/jobs/{jobid:int}',
        templateUrl: 'static/app/table/job.detail.html',
        controller: 'DetailCtrl'
      });
  });
