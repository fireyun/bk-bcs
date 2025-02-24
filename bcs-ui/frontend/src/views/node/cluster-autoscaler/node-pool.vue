<template>
  <BcsContent>
    <template #header>
      <HeaderNav :list="navList">
        <bcs-steps style="max-width: 360px" :steps="steps" :cur-step="curStep"></bcs-steps>
      </HeaderNav>
    </template>
    <div v-bkloading="{ isLoading }" class="node-pool">
      <keep-alive>
        <component
          :is="stepComMap[curStep]"
          :default-values="defaultValues"
          :schema="schema"
          :cluster="curCluster"
          :save-loading="saveLoading"
          v-if="!isLoading"
          @next="handleNextStep"
          @pre="handlePreStep"
          @confirm="handleConfirm"
        ></component>
      </keep-alive>
    </div>
  </BcsContent>
</template>
<script lang="ts">
import { defineComponent, ref, computed, onMounted, toRefs } from '@vue/composition-api';
import BcsContent from '../bcs-content.vue';
import HeaderNav from '../header-nav.vue';
import { useClusterList } from '@/views/cluster/use-cluster';
import $i18n from '@/i18n/i18n-setup';
import NodePoolInfo from './node-pool-info.vue';
import NodeConfig from './node-config.vue';
import $store from '@/store/index';
import Schema from '../resolve-schema';
import { mergeDeep } from '@/common/util';
import $router from '@/router/index';

export default defineComponent({
  components: {
    BcsContent,
    HeaderNav,
    NodePoolInfo,
    NodeConfig,
  },
  props: {
    clusterId: {
      type: String,
      default: '',
      required: true,
    },
  },
  setup(props, ctx) {
    const { clusterId } = toRefs(props);
    const { clusterList } = useClusterList(ctx);
    const curCluster = computed(() => ($store.state as any).cluster.clusterList
      ?.find(item => item.clusterID === clusterId.value) || {});
    const navList = computed(() => {
      const nav: any[] = [
        {
          title: clusterList.value.find(item => item.clusterID === clusterId.value)?.clusterName,
          link: {
            name: 'clusterDetail',
          },
        },
        {
          title: 'Cluster Autoscaler',
          link: {
            name: 'clusterDetail',
            query: {
              active: 'AutoScaler',
            },
          },
        },
        {
          title: $i18n.t('新建节点池'),
          link: null,
        },
      ];
      return nav;
    });
    const steps = ref([
      {
        title: $i18n.t('节点池信息'),
        icon: 1,
      },
      {
        title: $i18n.t('节点配置'),
        icon: 2,
      },
    ]);
    const curStep = ref(1);
    const curStepItem = computed<Record<string, any>>(() => steps.value
      .find((_, index) => index + 1 === curStep.value) || {});
    const stepComMap = {
      1: 'NodePoolInfo',
      2: 'NodeConfig',
    };
    const nodePoolData = ref<Record<string, any>>({});
    const handleNextStep = (data) => {
      nodePoolData.value = mergeDeep(nodePoolData.value, data);
      if (curStep.value + 1 <= steps.value.length) {
        curStep.value = curStep.value + 1;
      }
    };
    const handlePreStep = () => {
      curStep.value = curStep.value - 1;
    };

    const isLoading = ref(true);
    // 获取默认值
    const defaultValues = ref<any>(null);
    const schema = ref({});
    const handleGetCloudDefaultValues = async () => {
      const data = await $store.dispatch('clustermanager/resourceSchema', {
        $cloudID: 'selfProvisionCloud', // todo ieod暂时写死
        $name: 'nodegroup',
      });
      schema.value = data?.schema || {};
    };

    // 创建节点池
    const user = computed(() => $store.state.user);
    const saveLoading = ref(false);
    const handleConfirm = async () => {
      saveLoading.value = true;
      await handleCreateNodePool();
      saveLoading.value = false;
    };
    const handleCreateNodePool = async () => {
      const data = {
        ...nodePoolData.value,
        provider: 'selfProvisionCloud',
        clusterID: curCluster.value.clusterID,
        region: curCluster.value.region,
        creator: user.value.username,
      };
      console.log(data);
      const result = await $store.dispatch('clustermanager/createNodeGroup', data);
      if (result) {
        $router.push({
          name: 'clusterDetail',
          query: {
            active: 'AutoScaler',
          },
        });
      }
    };

    onMounted(async () => {
      isLoading.value = true;
      await handleGetCloudDefaultValues();
      defaultValues.value = Schema.getSchemaDefaultValue(schema.value);
      isLoading.value = false;
    });
    return {
      saveLoading,
      curCluster,
      isLoading,
      schema,
      defaultValues,
      navList,
      steps,
      curStep,
      curStepItem,
      stepComMap,
      handleNextStep,
      handlePreStep,
      handleConfirm,
    };
  },
});
</script>
<style lang="postcss" scoped>
.node-pool {
  margin: -24px;
  min-height: 200px;
}
</style>
