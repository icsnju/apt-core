'use strict';

angular.module('aptWebApp')
    .controller('SlaveCtrl', function($scope, $http) {
        $scope.nodes = [];
        $scope.searchKey = '';

        //update slaves status
        $scope.refresh = function() {
            $http.get('slave').then(response => {
                if (response) {
                    $scope.nodes = response.data;
                }
            });
        }

        $scope.refresh();
    })
    .controller('SlaveDetailCtrl', function($scope, $http, $stateParams) {
        $scope.slave = {};
        $scope.devices = [];
        $scope.tasks = [];
        $scope.ip = '';

        var id = $stateParams.id;
        $scope.ip = id;
        //console.log(JSON.stringify($stateParams));

        $scope.refresh = function() {
            $http.get('slave/' + id).then(response => {
                if (response) {
                    $scope.slave = response.data;
                    var taskmap = $scope.slave.TaskStates;
                    $scope.tasks = [];
                    for (var key in taskmap) {
                        $scope.tasks.push(taskmap[key]);
                    }
                    var devicemap = $scope.slave.DeviceStates;
                    $scope.devices = [];
                    for (var key in devicemap) {
                        $scope.devices.push(devicemap[key]);
                    }
                }
            });
        };

        $scope.getDevStatus = function(status) {
          if (status==0){
            return 'busy';
          }else{
            return 'free';
          }
        };

        //get status of the task
        $scope.getTaskStatus = function(status) {
            if (status == 0) {
                return 'waiting';
            } else if (status == 1) {
                return 'running';
            } else if (status == 2) {
                return 'finished';
            } else {
                return 'failed';
            }
        }

        //get the point color relay on the task status
        $scope.getPoColor = function(status) {
            if (status == 3) {
                return 'danger';
            } else if (status == 2) {
                return 'success';
            } else if (status == 1) {
                return 'info';
            } else {
                return 'wait';
            }
        }


        $scope.refresh();
    });
