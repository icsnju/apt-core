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
        $scope.ip='';

        var id = $stateParams.id;
        $scope.ip=id;
        console.log(JSON.stringify($stateParams));

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
        }

        $scope.refresh();
    });
