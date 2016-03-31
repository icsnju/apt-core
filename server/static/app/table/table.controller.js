'use strict';

angular.module('aptWebApp')
  .controller('JobCtrl', function($scope, $http, $state) {
    $scope.jobs = [];

    $http.get('job').then(response => {
      if (response) {
        $scope.jobs = response.data;
      }
    });

  })
  .controller('DetailCtrl', function($scope, $http, $state, $stateParams, $window) {
    $scope.job = {};
    $scope.tasks = [];
    var jid = $stateParams.jobid
    $http.get('job/' + jid).then(response => {
      if (response) {
        $scope.job = response.data;
        var taskmap = $scope.job.taskmap;
        for (var key in taskmap) {
          $scope.tasks.push(taskmap[key]);
        }
      }
    });

    $scope.DownloadResult = function(did, index) {
      if ($scope.tasks[index].state != 2) {
        return;
      }
      var dlurl = 'download/task?deviceid=' + did + '&' + 'jobid=' + jid;
      $window.open(dlurl);
    }

  });
