<template>
  <keep-alive>
    <component
      :is="componentName"
      :cur-project="curProject"
      :category="category"
      :namespace="namespace"
      :name="name"
      :kind="kind"
      :crd="crd"
      :hidden-operate="true"
    ></component>
  </keep-alive>
  <!-- <router-view v-else :key="$route.path"></router-view> -->
</template>

<script>
const deployments = () => import(
  /* webpackChunkName: 'app-list' */'./k8s/deployments');
const deploymentsInstanceDetail = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const deploymentsInstanceDetail2 = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const deploymentsContainerDetail = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/deployments-container');
const deploymentsContainerDetail2 = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/deployments-container2');
const deploymentsInstantiation = () => import(
  /* webpackChunkName: 'app-instantiation' */'./k8s/deployments-instantiation');
const daemonset = () => import(
  /* webpackChunkName: 'app-list' */'./k8s/daemonset');
const daemonsetInstanceDetail = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const daemonsetInstanceDetail2 = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const daemonsetContainerDetail = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/daemonset-container');
const daemonsetContainerDetail2 = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/daemonset-container2');
const daemonsetInstantiation = () => import(
  /* webpackChunkName: 'app-instantiation' */'./k8s/daemonset-instantiation');

const job = () => import(
  /* webpackChunkName: 'app-list' */'./k8s/job');
const jobInstanceDetail = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const jobInstanceDetail2 = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const jobContainerDetail = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/job-container');
const jobContainerDetail2 = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/job-container2');
const jobInstantiation = () => import(
  /* webpackChunkName: 'app-instantiation' */'./k8s/job-instantiation');

const statefulset = () => import(
  /* webpackChunkName: 'app-list' */'./k8s/statefulset');
const statefulsetInstanceDetail = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const statefulsetInstanceDetail2 = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const statefulsetContainerDetail = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/statefulset-container');
const statefulsetContainerDetail2 = () => import(
  /* webpackChunkName: 'app-container' */'./k8s/statefulset-container2');
const statefulsetInstantiation = () => import(
  /* webpackChunkName: 'app-instantiation' */'./k8s/statefulset-instantiation');
const gamestatefulset = () => import(
  /* webpackChunkName: 'app-list' */'./k8s/gamestatefulset');
const gamestatefulSetsInstanceDetail = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const gamedeployments = () => import(
  /* webpackChunkName: 'app-list' */'./k8s/gamedeployments');
const gamedeploymentsInstanceDetail = () => import(
  /* webpackChunkName: 'app-instance' */'@/views/dashboard/workload/detail/index.vue');
const customobjects = () => import(
  /* webpackChunkName: 'app-list' */'./k8s/customobjects');

export default {
  name: 'DetailIndex',
  components: {

    deployments,
    deploymentsInstanceDetail,
    deploymentsInstanceDetail2,
    deploymentsContainerDetail,
    deploymentsContainerDetail2,
    deploymentsInstantiation,

    daemonset,
    daemonsetInstanceDetail,
    daemonsetInstanceDetail2,
    daemonsetContainerDetail,
    daemonsetContainerDetail2,
    daemonsetInstantiation,

    job,
    jobInstanceDetail,
    jobInstanceDetail2,
    jobContainerDetail,
    jobContainerDetail2,
    jobInstantiation,

    statefulset,
    statefulsetInstanceDetail,
    statefulsetInstanceDetail2,
    statefulsetContainerDetail,
    statefulsetContainerDetail2,
    statefulsetInstantiation,

    gamestatefulset,
    gamestatefulSetsInstanceDetail,
    gamedeployments,
    gamedeploymentsInstanceDetail,
    customobjects,
  },
  data() {
    return {
      k8sPathNameList: [
        'deployments',
        'deploymentsInstanceDetail',
        'deploymentsInstanceDetail2',
        'deploymentsContainerDetail',
        'deploymentsContainerDetail2',
        'deploymentsInstantiation',
        'daemonset',
        'daemonsetInstanceDetail',
        'daemonsetInstanceDetail2',
        'daemonsetContainerDetail',
        'daemonsetContainerDetail2',
        'daemonsetInstantiation',
        'job',
        'jobInstanceDetail',
        'jobInstanceDetail2',
        'jobContainerDetail',
        'jobContainerDetail2',
        'jobInstantiation',
        'statefulset',
        'statefulsetInstanceDetail',
        'statefulsetInstanceDetail2',
        'statefulsetContainerDetail',
        'statefulsetContainerDetail2',
        'statefulsetInstantiation',

        'gamestatefulset',
        'gamestatefulSetsInstanceDetail',
        'gamedeployments',
        'gamedeploymentsInstanceDetail',
        'customobjects',
      ],
      componentName: '',
      projectId: '',
      projectCode: '',
      curProject: null,
    };
  },
  computed: {
    onlineProjectList() {
      return this.$store.state.sideMenu.onlineProjectList;
    },
    curProjectCode() {
      return this.$store.state.curProjectCode;
    },
    curProjectId() {
      return this.$store.state.curProjectId;
    },
    category() {
      const categoryMap = {
        deploymentsInstanceDetail: 'deployments',
        deploymentsInstanceDetail2: 'deployments',
        daemonsetInstanceDetail: 'daemonsets',
        daemonsetInstanceDetail2: 'daemonsets',
        statefulsetInstanceDetail: 'statefulsets',
        statefulsetInstanceDetail2: 'statefulsets',
        jobInstanceDetail: 'jobs',
        jobInstanceDetail2: 'jobs',
        gamedeploymentsInstanceDetail: 'custom_objects',
        gamestatefulSetsInstanceDetail: 'custom_objects',
      };
      return categoryMap[this.$route.name];
    },
    namespace() {
      return this.$route.query?.namespace;
    },
    name() {
      return this.$route.query?.name;
    },
    kind() {
      const kindMap = {
        deploymentsInstanceDetail: 'Deployment',
        deploymentsInstanceDetail2: 'Deployment',
        daemonsetInstanceDetail: 'DaemonSet',
        daemonsetInstanceDetail2: 'DaemonSet',
        statefulsetInstanceDetail: 'StatefulSet',
        statefulsetInstanceDetail2: 'StatefulSet',
        jobInstanceDetail: 'Job',
        jobInstanceDetail2: 'Job',
        gamedeploymentsInstanceDetail: 'GameDeployment',
        gamestatefulSetsInstanceDetail: 'GameStatefulSet',
      };
      return kindMap[this.$route.name];
    },
    crd() {
      const crdMap = {
        gamedeploymentsInstanceDetail: 'gamedeployments.tkex.tencent.com',
        gamestatefulSetsInstanceDetail: 'gamestatefulsets.tkex.tencent.com',
      };
      return crdMap[this.$route.name];
    },
  },
  mounted() {
    this.setCurProject();
  },
  methods: {
    /**
             * 设置 curProject
             */
    setCurProject() {
      const len = this.onlineProjectList.length;
      if (len) {
        const firstProject = this.onlineProjectList[0];
        this.projectId = this.$route.params.projectId
                        || this.curProjectId
                        || firstProject.project_id;

        this.projectCode = this.$route.params.projectCode
                        || this.curProjectCode
                        || firstProject.english_name;

        this.curProject = this.onlineProjectList.find(p => p.project_id === this.projectId);
        if (this.curProject) {
          this.setComponent();
        }
      }
    },

    /**
             * 设置动态组件
             */
    setComponent() {
      const routeName = this.$route.name;
      if (this.k8sPathNameList.indexOf(routeName) <= -1) {
        this.componentName = 'deployments';
      } else {
        this.componentName = routeName;
      }
    },
  },
};
</script>
