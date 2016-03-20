'use strict';

angular.module('aptWebApp')
  .controller('SubmitCtrl', function($scope, $http, Upload, $state) {

    //init
    $scope.devices = [];
    $scope.monkey = {};

    $scope.frameErr = false;
    $scope.selectorErr = false;
    var idList = [];


    //get devices from server
//    $http.get('/api/devices').then(response => {
//      if (response) {
//        $scope.devices = response.data;
//        for (var i = 0; i < $scope.devices.length; i++) {
//          $scope.devices[i].check = false;
//        }
//      }
//    });

    $scope.checkEmpty = function() {

      //get all selected devices
      var devices = $scope.devices;
      idList = [];
      for (var i = 0; i < devices.length; i++) {
        if (devices[i].check) {
          idList.push(devices[i].Id);
        }
      }

      $scope.frameErr = false;
      $scope.selectorErr = false;
      //check empty input
      if (!$scope.monkey.file || !$scope.monkey.pkg || !$scope.monkey.arg) {
        $scope.frameErr = true;
      } else if (idList.length <= 0) {
        $scope.selectorErr = true;
      }
    }

    //submit this job
    $scope.submitJob = function() {

      $scope.checkEmpty();
      if ($scope.frameErr || $scope.selectorErr) {
        return;
      }

      //create SubJob struct
      var SubJob = {};
      SubJob.FrameKind = 'monkey';
      SubJob.Frame = {};
      SubJob.Frame.AppPath = $scope.monkey.file.name;
      SubJob.Frame.PkgName = $scope.monkey.pkg;
      SubJob.Frame.Argu = $scope.monkey.arg;
      SubJob.FilterKind = 'specify_devices';
      SubJob.Filter = {};
      SubJob.Filter.IdList = idList;

      //submit requirement and files to server
      Upload.upload({
        url: '/api/jobs',
        method: 'POST',
        data: {
          file: $scope.monkey.file,
          job: SubJob
        }
      }).then(function(resp) {
        console.log('Success ' + resp.config.data.file.name + 'uploaded. Response: ' + resp.data);
        $state.go('table');
      }, function(resp) {
        console.log('Error status: ' + resp.status);
      }, function(evt) {
        var progressPercentage = parseInt(100.0 * evt.loaded / evt.total);
        console.log('progress: ' + progressPercentage + '% ' + evt.config.data.file.name);
      });

    }

  });
