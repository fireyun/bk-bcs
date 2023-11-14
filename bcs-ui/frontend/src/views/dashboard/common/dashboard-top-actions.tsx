import { defineComponent, computed } from '@vue/composition-api';

import './dashboard-top-actions.css';

export default defineComponent({
  name: 'DashboardTopActions',
  setup(props, ctx) {
    const { $router, $route } = ctx.root;

    const projectId = computed(() => $route.params.projectId).value;
    const projectCode = computed(() => $route.params.projectCode).value;

    const goBCS = () => {
      $router.push({
        name: 'clusterMain',
        params: {
          projectId,
          projectCode,
        },
      });
    };

    return {
      goBCS,
    };
  },
  render() {
    return (
            <div class="dashboard-top-actions">
            </div>
  );
  },
});
