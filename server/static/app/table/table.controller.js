'use strict';

angular.module('aptWebApp')
  .controller('TableCtrl', function($scope, $http, $state) {
    $scope.jt = {};
    $scope.jt.bsTableControl = {};
    $scope.jobs = [];

    $http.get('job').then(response => {
      if (response) {
        $scope.jobs = response.data;
      }
    });

  })
  .controller('TasksCtrl', function($scope, $http, $state, $stateParams) {
    $scope.job={};
    $scope.tasks=[];
    $http.get('job/'+$stateParams.jobid).then(response => {
      if (response) {
        $scope.job = response.data;
        var taskmap=$scope.job.taskmap;
        for(var key in taskmap){
          $scope.tasks.push(taskmap[key]);
        }
      }
    });
  });
