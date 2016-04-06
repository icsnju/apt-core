'use strict';

angular.module('aptWebApp')
    .controller('JobCtrl', function($scope, $http, $state) {
        $scope.jobs = [];
        $scope.searchKey = '';
        $scope.statusKey = {};

        $scope.setStatusKey = function(key) {
            // console.log(key);
            var statusKey;
            if (key == 'finish') {
                statusKey = {
                    status: '100'
                };
            } else if (key == 'fail') {
                statusKey = {
                    status: '-1'
                };
            } else if (key == 'run') {
                statusKey = function(job) {
                    return job.status < 100 && job.status != -1;
                }
            } else {
                statusKey = {};
            }

            $scope.statusKey = statusKey;
        }

        $scope.getPercent = function(status) {
            if (status == -1) {
                return '100%';
            } else {
                return status + '%';
            }
        }

        $scope.getProColor = function(status) {
            if (status < 0) {
                return 'danger';
            } else if (status == 100) {
                return 'success';
            } else {
                return 'info';
            }
        }

        $http.get('job').then(response => {
            if (response) {
                $scope.jobs = response.data;
            }
        });

    })
    .controller('DetailCtrl', function($scope, $http, $state, $stateParams, $window) {
        $scope.job = {};
        $scope.tasks = [];
        $scope.searchKey = '';
        $scope.statusKey = {};


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

        //download result file frome server
        $scope.downloadResult = function(did, index) {
            if ($scope.tasks[index].state != 2) {
                return;
            }
            var dlurl = 'download/task?deviceid=' + did + '&' + 'jobid=' + jid;
            $window.open(dlurl);
        }

        //set status search key
        $scope.setStatusKey = function(key) {
            // console.log(key);
            var statusKey;
            if (key == 'finish') {
                statusKey = {
                    state: 2
                };
            } else if (key == 'fail') {
                statusKey = {
                    state: 3
                };
            } else if (key == 'run') {
                statusKey = function(task) {
                    return task.state == 0 || task.state == 1;
                }
            } else {
                statusKey = {};
            }

            $scope.statusKey = statusKey;
        }

        //get status of the task
        $scope.getStatus = function(status) {
            if (status == 0 || status == 1) {
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
            } else {
                return 'info';
            }
        }

    });
