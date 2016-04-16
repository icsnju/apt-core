'use strict';

angular.module('aptWebApp')
    .controller('JobCtrl', function($scope, $http, $window) {
        $scope.jobs = [];
        $scope.searchKey = '';
        $scope.statusKey = {};

        //set status filter key
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

        //get progress
        $scope.getPercent = function(status) {
            if (status == -1) {
                return '100%';
            } else {
                return status + '%';
            }
        }

        //get progress bar color
        $scope.getProColor = function(status) {
            if (status < 0) {
                return 'danger';
            } else if (status == 100) {
                return 'success';
            } else {
                return 'info';
            }
        }

        //if this option clickable
        $scope.getClickAble = function(status, type) {
            if (type == 'kill') {
                if (status >= 0 && status < 100) {
                    return 'options';
                } else {
                    return 'nooptions'
                }
            } else {
                if (status >= 0 && status < 100) {
                    return 'nooptions';
                } else {
                    return 'options'
                }
            }
        }

        //kill a job
        $scope.killJob = function(id, status) {
            if (status < 0 || status >= 100) {
                return;
            }

            $http.put('job/' + id)
                .then(
                    function(res) {
                        $http.get('job').then(response => {
                            if (response) {
                                $scope.jobs = response.data;
                            }
                        });
                    },
                    function(res) {
                        console.log('Error status: ' + res.status);
                    }
                );
        }

        //delete a job
        $scope.deleteJob = function(id, status) {
            if (status >= 0 && status < 100) {
                return;
            }

            $http.delete('job/' + id)
                .then(
                    function(res) {
                        $http.get('job').then(response => {
                            if (response) {
                                $scope.jobs = response.data;
                            }
                        });
                    },
                    function(res) {
                        console.log('Error status: ' + res.status);
                    }
                );
        }

        //download testing result of this job
        $scope.downloadJob = function(id, status) {
            if (status >= 0 && status < 100) {
                return;
            }
            var dlurl = 'download/job?jobid=' + id;
            $window.open(dlurl);
        }

        //update jobs status
        $scope.refresh = function() {
            $http.get('job').then(response => {
                if (response) {
                    $scope.jobs = response.data;
                }
            });
        }

        $scope.refresh();
    })
    .controller('DetailCtrl', function($scope, $http, $stateParams, $window) {
        $scope.job = {};
        $scope.tasks = [];
        $scope.searchKey = '';
        $scope.statusKey = {};


        var jid = $stateParams.jobid;

        $scope.refresh = function() {
            $http.get('job/' + jid).then(response => {
                if (response) {
                    $scope.job = response.data;
                    var taskmap = $scope.job.taskmap;
                    $scope.tasks = [];
                    for (var key in taskmap) {
                        $scope.tasks.push(taskmap[key]);
                    }
                }
            });
        }

        $scope.refresh();

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

        //get the point color relay on the task status
        $scope.getDownAble = function(status) {
            if (status == 2) {
                return 'options';
            } else {
                return 'nooptions';
            }
        }

    });
